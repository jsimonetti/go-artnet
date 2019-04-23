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

// we poll for new nodes every 3 seconds
var pollInterval = 3 * time.Second

// ControlledNode hols the configuration of a node we control
type ControlledNode struct {
	LastSeen   time.Time
	Node       NodeConfig
	UDPAddress net.UDPAddr

	Sequence   uint8
	DMXBuffer  map[Address][512]byte
	LastUpdate time.Time
	nodeLock   sync.Mutex
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

	pollTicker *time.Ticker
	gcTicker   *time.Ticker
}

// NewController return a Controller
func NewController(name string, ip net.IP, log Logger) *Controller {
	return &Controller{
		cNode: NewNode(name, code.StController, ip, log),
		log:   log,
	}
}

// Start will start this controller
func (c *Controller) Start() error {
	c.OutputAddress = make(map[Address]*ControlledNode)
	c.InputAddress = make(map[Address]*ControlledNode)
	c.shutdownCh = make(chan struct{})
	c.cNode.log = c.log.With(Fields{"type": "Node"})
	c.log = c.log.With(Fields{"type": "Controller"})
	if err := c.cNode.Start(); err != nil {
		return fmt.Errorf("failed to start controller node: %v", err)
	}

	c.pollTicker = time.NewTicker(pollInterval)
	c.gcTicker = time.NewTicker(pollInterval)

	go c.pollLoop()
	go c.dmxUpdateLoop()
	return c.cNode.shutdownErr
}

// Stop will stop this controller
func (c *Controller) Stop() {
	c.pollTicker.Stop()
	c.gcTicker.Stop()
	c.cNode.Stop()

	select {
	case <-c.cNode.shutdownCh:
	}

	close(c.shutdownCh)
}

// pollLoop will routinely poll for new nodes
func (c *Controller) pollLoop() {
	artPoll := &packet.ArtPollPacket{
		TalkToMe: new(code.TalkToMe).WithReplyOnChange(true),
		Priority: code.DpAll,
	}

	// create an ArtPoll packet to send out periodically
	b, err := artPoll.MarshalBinary()
	if err != nil {
		c.log.With(Fields{"err": err}).Error("error creating ArtPoll packet")
		return
	}

	// send ArtPollPacket
	c.cNode.sendCh <- netPayload{
		address: broadcastAddr,
		data:    b,
	}
	c.cNode.pollCh <- packet.ArtPollPacket{}

	// loop until shutdown
	for {
		select {
		case <-c.pollTicker.C:
			// send ArtPollPacket
			c.cNode.sendCh <- netPayload{
				address: broadcastAddr,
				data:    b,
			}
			c.cNode.pollCh <- packet.ArtPollPacket{}

		case <-c.gcTicker.C:
			// clean up old nodes
			c.gcNode()

		case p := <-c.cNode.pollReplyCh:
			cfg := ConfigFromArtPollReply(p)
			if cfg.Type != code.StNode {
				// we don't care for ArtNet devices other then nodes for now @todo
				continue
			}

			if err := c.updateNode(cfg); err != nil {
				c.log.With(Fields{"err": err}).Error("error updating node")
			}

		case <-c.shutdownCh:
			return
		}
	}
}

// SendDMXToAddress will set the DMXBuffer for a destination address
// and update the node
func (c *Controller) SendDMXToAddress(dmx [512]byte, address Address) {
	c.log.With(Fields{"address": address.String()}).Debug("received update channels")

	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	var cn *ControlledNode
	var ok bool

	if cn, ok = c.OutputAddress[address]; !ok {
		c.log.With(Fields{"address": address.String()}).Error("could not find node for address")
		return
	}
	err := cn.setDMXBuffer(dmx, address)
	if err != nil {
		c.log.With(Fields{"err": err, "address": address.String()}).Error("error setting buffer on address")
		return
	}

	// get an ArtDMXPacket for this node
	b, err := cn.dmxUpdate(address)
	if err != nil {
		c.log.With(Fields{"err": err}).Error("error getting packet for dmxUpdate")
		return
	}

	cn.LastUpdate = time.Now()

	c.cNode.sendCh <- netPayload{
		address: cn.UDPAddress,
		data:    b,
	}

}

// dmxUpdateLoop will periodically update nodes until shutdown
func (c *Controller) dmxUpdateLoop() {
	// we force update nodes every 15 seconds
	updateTicker := time.NewTicker(250 * time.Millisecond)
	updateAfter := 800 * time.Millisecond

	// loop until shutdown
	for {
		select {
		case <-updateTicker.C:

			now := time.Now()
			// send DMX buffer update
			c.nodeLock.Lock()
			for address, node := range c.OutputAddress {
				// only update if it has been X seconds
				if node.LastUpdate.Add(updateAfter).Before(now) {
					// get an ArtDMXPacket for this node
					b, err := node.dmxUpdate(address)
					if err != nil {
						c.log.With(Fields{"err": err, "address": address.String()}).Error("error getting buffer for address")
						break
					}
					node.LastUpdate = now
					c.cNode.sendCh <- netPayload{
						address: node.UDPAddress,
						data:    b,
					}
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
			c.log.With(Fields{"node": cfg.Name, "ip": cfg.IP.String()}).Debug("updated node")
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
	c.log.With(Fields{"node": cfg.Name, "ip": cfg.IP.String()}).Debug("added node")
	node := &ControlledNode{
		Node:       cfg,
		DMXBuffer:  buf,
		LastSeen:   time.Now(),
		Sequence:   0,
		UDPAddress: net.UDPAddr{IP: cfg.IP, Port: int(packet.ArtNetPort)},
		LastUpdate: time.Now(),
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

	// nodes are stale after 5 missed ArtPoll's
	//staleAfter, _ := time.ParseDuration(fmt.Sprintf("%ds", 5*pollInterval))
	staleAfter := 7 * time.Second

start:
	for i := range c.Nodes {
		if c.Nodes[i].LastSeen.Add(staleAfter).Before(time.Now()) {
			// it has been more then X seconds since we saw this node. remove it now.
			c.log.With(Fields{"node": c.Nodes[i].Node.Name, "ip": c.Nodes[i].Node.IP.String()}).Debug("remove stale node")

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
