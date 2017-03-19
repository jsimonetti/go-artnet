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

// ControlledNode hols the configuration of a node we control
type ControlledNode struct {
	LastSeen time.Time
	Node     Config
}

// Controller holds the information for a controller
type Controller struct {
	// cNode is the Node for the cNode
	cNode *Node

	// Nodes is a slice of nodes that are seen by this controller
	Nodes         []ControlledNode
	OutputAddress map[Address]*ControlledNode
	InputAddress  map[Address]*ControlledNode
	nodeLock      sync.Mutex

	shutdownCh chan struct{}
}

// NewController return a Controller
func NewController(name string, ip net.IP) *Controller {
	return &Controller{
		cNode: NewNode(name, code.StController, ip),
	}
}

// Start will start this controller
func (c *Controller) Start() error {
	c.OutputAddress = make(map[Address]*ControlledNode)
	c.InputAddress = make(map[Address]*ControlledNode)
	go c.pollLoop()
	return c.cNode.Start()
}

// Stop will stop this controller
func (c *Controller) Stop() {
	c.cNode.Stop()
	close(c.shutdownCh)
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
		fmt.Printf("error creating ArtPoll packet: %s\n", err)
		return
	}

	// create an ArtPollReply packet to send out with the ArtPoll packet
	me, err := new(packet.ArtPollReplyPacket).MarshalBinary()
	if err != nil {
		fmt.Printf("error creating ArtPollReply packet for self: %s\n", err)
		return
	}

	// loop untill shutdown
	for {
		select {
		case <-pollTicker.C:
			// send ArtPollPacket
			c.cNode.sendCh <- &netPayload{data: b}

			// we should always reply to our own polls to let other controllers know we are here
			c.cNode.sendCh <- &netPayload{data: me}

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

// updateNode will add a Node to the list of known nodes
// this assumes that there are no universe address collisions
// in the future we should probably be prepared to handle that too
func (c *Controller) updateNode(cfg Config) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for i, n := range c.Nodes {
		if bytes.Equal(cfg.IP, n.Node.IP) {
			// update this node, since we allready know about it
			fmt.Printf("updated node: %s, %s\n", cfg.Name, cfg.IP.String())
			// remove references to this node from the output map
			for _, port := range c.Nodes[i].Node.OutputPorts {
				delete(c.OutputAddress, port.Address)
			}
			for _, port := range c.Nodes[i].Node.InputPorts {
				delete(c.InputAddress, port.Address)
			}
			c.Nodes[i].Node = cfg
			c.Nodes[i].LastSeen = time.Now()
			// add references to this node from the output map
			for _, port := range c.Nodes[i].Node.OutputPorts {
				c.OutputAddress[port.Address] = &c.Nodes[i]
			}
			for _, port := range c.Nodes[i].Node.InputPorts {
				c.InputAddress[port.Address] = &c.Nodes[i]
			}
			return nil
		}
	}
	// new node, add it to our known nodes
	fmt.Printf("added node: %s, %s\n", cfg.Name, cfg.IP.String())
	c.Nodes = append(c.Nodes, ControlledNode{Node: cfg, LastSeen: time.Now()})

	return nil
}

// deleteNode will delete a Node from the list of known nodes
func (c *Controller) deleteNode(node Config) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for i, n := range c.Nodes {
		if bytes.Equal(node.IP, n.Node.IP) {
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
			fmt.Printf("remove stale node: %s, %s\n", c.Nodes[i].Node.Name, c.Nodes[i].Node.IP.String())
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
