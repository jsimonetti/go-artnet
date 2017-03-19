package artnet

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

var broadcastAddr = net.UDPAddr{
	IP:   []byte{0x02, 0xff, 0xff, 0xff},
	Port: int(packet.ArtNetPort),
}

// ControlledNode hols the configuration of a node we control
type ControlledNode struct {
	LastSeen   time.Time
	Node       NodeConfig
	UDPAddress net.UDPAddr

	Sequence  uint8
	DMXBuffer map[Address][512]byte
	nodeLock  sync.Mutex
}

// setDMXBuffer will update the buffer on a universe address
func (cn *ControlledNode) setDMXBuffer(dmx [512]byte, address Address) error {
	cn.nodeLock.Lock()
	defer cn.nodeLock.Unlock()

	if _, ok := cn.DMXBuffer[address]; !ok {
		return fmt.Errorf("unknown address for controlled node")
	}
	cn.DMXBuffer[address] = dmx
	return nil
}

// dmxUpdate will create an ArtDMXPacket and marshal it into bytes
func (cn *ControlledNode) dmxUpdate(address Address) (b []byte, err error) {
	var buf [512]byte
	var ok bool

	cn.nodeLock.Lock()
	defer cn.nodeLock.Unlock()

	if buf, ok = cn.DMXBuffer[address]; !ok {
		return nil, fmt.Errorf("unknown address for controlled node")
	}

	cn.Sequence++
	p := &packet.ArtDMXPacket{
		Sequence: cn.Sequence,
		SubUni:   address.SubUni,
		Net:      address.Net,
		Length:   uint16(len(cn.DMXBuffer)),
		Data:     buf,
	}
	b, err = p.MarshalBinary()
	return
}

// Controller holds the information for a controller
type Controller struct {
	// cNode is the Node for the cNode
	cNode *Node

	// Nodes is a slice of nodes that are seen by this controller
	Nodes         []*ControlledNode
	OutputAddress map[Address]*ControlledNode
	InputAddress  map[Address]*ControlledNode
	nodeLock      sync.Mutex

	shutdownCh chan struct{}

	log Logger
}

// NewController return a Controller
func NewController(name string, ip net.IP) *Controller {
	return &Controller{
		cNode: NewNode(name, code.StController, ip),
		log:   NewLogger(),
	}
}

// Start will start this controller
func (c *Controller) Start() error {
	c.OutputAddress = make(map[Address]*ControlledNode)
	c.InputAddress = make(map[Address]*ControlledNode)
	c.shutdownCh = make(chan struct{})
	c.cNode.Start()
	go c.pollLoop()
	go c.dmxUpdateLoop()
	return c.cNode.shutdownErr
}

// Stop will stop this controller
func (c *Controller) Stop() {
	close(c.shutdownCh)
	c.cNode.Stop()
}

// pollLoop will routinely poll for new nodes
func (c *Controller) pollLoop() {
	// we poll for new nodes every 5 seconds
	pollTicker := time.NewTicker(5 * time.Second)

	// we garbagecollect every 30 seconds
	gcTicker := time.NewTicker(30 * time.Second)

	artPoll := &packet.ArtPollPacket{
		TalkToMe: new(code.TalkToMe).WithReplyOnChange(true),
		Priority: code.DpAll,
	}

	// create an ArtPoll packet to send out periodically
	b, err := artPoll.MarshalBinary()
	if err != nil {
		c.log.With(Fields{"err": err}).Printf("error creating ArtPoll packet")
		return
	}

	// create an ArtPollReply packet to send out with the ArtPoll packet
	me, err := new(packet.ArtPollReplyPacket).MarshalBinary()
	if err != nil {
		c.log.With(Fields{"err": err}).Printf("error creating ArtPollReply packet for self")
		return
	}

	// loop untill shutdown
	for {
		select {
		case <-pollTicker.C:
			// send ArtPollPacket
			c.cNode.sendCh <- netPayload{
				address: broadcastAddr,
				data:    b,
			}

			// we should always reply to our own polls to let other controllers know we are here
			c.cNode.sendCh <- netPayload{
				address: broadcastAddr,
				data:    me,
			}

		case <-gcTicker.C:
			// clean up old nodes
			c.gcNode()

		case p := <-c.cNode.pollReplyCh:
			cfg := ConfigFromArtPollReply(p)
			c.updateNode(cfg)

		case <-c.shutdownCh:
			return
		}
	}
}

// SendDMXToAddress will set the DMXBuffer for a destination address
// and update the node
func (c *Controller) SendDMXToAddress(dmx [512]byte, address Address) {
	fmt.Printf("received update channels to %x, %x, %x\n", dmx[0], dmx[1], dmx[2])

	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	var cn *ControlledNode
	var ok bool

	if cn, ok = c.OutputAddress[address]; !ok {
		c.log.With(Fields{"address": address.String()}).Printf("could not find node for address")
		return
	}
	err := cn.setDMXBuffer(dmx, address)
	if err != nil {
		c.log.With(Fields{"err": err, "address": address.String()}).Printf("error setting buffer on address")
		return
	}

	// get an ArtDMXPacket for this node
	b, err := cn.dmxUpdate(address)
	if err != nil {
		c.log.With(Fields{"err": err}).Printf("error getting packet for dmxUpdate")
		return
	}

	c.cNode.sendCh <- netPayload{
		address: cn.UDPAddress,
		//address: broadcastAddr,
		data: b,
	}

}

// dmxUpdateLoop will periodically update nodes until shutdown
func (c *Controller) dmxUpdateLoop() {
	// we force update nodes every 14 seconds
	updateTicker := time.NewTicker(14 * time.Second)

	// loop untill shutdown
	for {
		select {
		case <-updateTicker.C:
			// send DMX buffer update
			c.nodeLock.Lock()
			for address, node := range c.OutputAddress {
				// get an ArtDMXPacket for this node
				b, err := node.dmxUpdate(address)
				if err != nil {
					c.log.With(Fields{"err": err, "address": address.String()}).Printf("error getting buffer for address")
					break
				}
				c.cNode.sendCh <- netPayload{
					address: node.UDPAddress,
					//address: broadcastAddr,
					data: b,
				}
			}
			c.nodeLock.Unlock()

		case <-c.shutdownCh:
			return
		}
	}
}

// updateNode will add a Node to the list of known nodes
// this assumes that there are no universe address collisions
// in the future we should probably be prepared to handle that too
func (c *Controller) updateNode(cfg NodeConfig) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for i := range c.Nodes {
		if bytes.Equal(cfg.IP, c.Nodes[i].Node.IP) {
			// update this node, since we allready know about it
			c.log.With(Fields{"node": cfg.Name, "ip": cfg.IP.String()}).Printf("updated node")
			// remove references to this node from the output map
			for _, port := range c.Nodes[i].Node.OutputPorts {
				delete(c.OutputAddress, port.Address)
			}
			for _, port := range c.Nodes[i].Node.InputPorts {
				delete(c.InputAddress, port.Address)
			}
			c.Nodes[i].Node = cfg
			c.Nodes[i].LastSeen = time.Now()
			// add references to this node to the output map
			for _, port := range c.Nodes[i].Node.OutputPorts {
				c.OutputAddress[port.Address] = c.Nodes[i]
			}
			for _, port := range c.Nodes[i].Node.InputPorts {
				c.InputAddress[port.Address] = c.Nodes[i]
			}
			return nil
		}
	}

	// create an empty DMX buffer. This will blackout the node entirely
	buf := make(map[Address][512]byte)
	for _, port := range cfg.OutputPorts {
		buf[port.Address] = [512]byte{}
	}

	// new node, add it to our known nodes
	c.log.With(Fields{"node": cfg.Name, "ip": cfg.IP.String()}).Printf("added node")
	node := &ControlledNode{
		Node:       cfg,
		DMXBuffer:  buf,
		LastSeen:   time.Now(),
		Sequence:   0,
		UDPAddress: net.UDPAddr{IP: cfg.IP, Port: int(packet.ArtNetPort)},
	}
	c.Nodes = append(c.Nodes, node)

	// add references to this node to the output map
	for _, port := range node.Node.OutputPorts {
		c.OutputAddress[port.Address] = node
	}
	for _, port := range node.Node.InputPorts {
		c.InputAddress[port.Address] = node
	}

	return nil
}

// deleteNode will delete a Node from the list of known nodes
func (c *Controller) deleteNode(node NodeConfig) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for i := range c.Nodes {
		if bytes.Equal(node.IP, c.Nodes[i].Node.IP) {
			// node found, remove it from the list
			// remove references to this node from the output map
			for _, port := range c.Nodes[i].Node.OutputPorts {
				delete(c.OutputAddress, port.Address)
			}
			for _, port := range c.Nodes[i].Node.InputPorts {
				delete(c.InputAddress, port.Address)
			}
			c.Nodes = append(c.Nodes[:i], c.Nodes[i+1:]...)
		}
	}

	return fmt.Errorf("no known node with this ip known, ip: %s", node.IP)
}

// gcNode will remove stale Nodes from the list of known nodes
// it will loop through the list of nodes and remove nodes older then X seconds
func (c *Controller) gcNode() {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	// we use X = 10 here, configurable in the future
	staleAfter := 10 * time.Second

start:
	for i := range c.Nodes {
		if c.Nodes[i].LastSeen.Add(staleAfter).Before(time.Now()) {
			// it has been more then X seconds since we saw this node. remove it now.
			c.log.With(Fields{"node": c.Nodes[i].Node.Name, "ip": c.Nodes[i].Node.IP.String()}).Printf("remove stale node")

			// remove references to this node from the output map
			for _, port := range c.Nodes[i].Node.OutputPorts {
				delete(c.OutputAddress, port.Address)
			}
			for _, port := range c.Nodes[i].Node.InputPorts {
				delete(c.InputAddress, port.Address)
			}
			// remove node
			c.Nodes = append(c.Nodes[:i], c.Nodes[i+1:]...)
			goto start
		}
	}
}
