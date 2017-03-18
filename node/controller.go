package node

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

// controlNode hols a node configuration
type controlNode struct {
	lastSeen time.Time
	node     Config
}

// Controller holds the information for a controller
type Controller struct {
	// Node is the controller itself
	Node

	// Nodes is a slice of nodes that are seen by this controller
	Nodes    []controlNode
	nodeLock sync.Mutex

	shutdownCh chan struct{}
}

// Start will start this controller
func (c *Controller) Start() error {
	go c.pollLoop()
	return c.Node.Start()
}

// Stop will stop this controller
func (c *Controller) Stop() {
	c.Node.Stop()
	close(c.shutdownCh)
}

// pollLoop will routinely poll for new nodes
func (c *Controller) pollLoop() {
	// we poll for new nodes every 3 seconds
	timer := time.NewTicker(3 * time.Second)
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
		case <-timer.C:
			// send ArtPollPacket
			c.Node.sendCh <- &netPayload{data: b}

			// we should always reply to our own polls to let other controllers know we are here
			c.Node.sendCh <- &netPayload{data: me}

			// clean up old nodes
			go c.gcNode()

		case p := <-c.Node.pollReplyCh:
			cfg := ConfigFromArtPollReply(p)
			c.updateNode(cfg)

		case <-c.shutdownCh:
			return
		}
	}
}

// updateNode will add a Node to the list of known nodes
func (c *Controller) updateNode(cfg Config) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for i, n := range c.Nodes {
		if bytes.Equal(cfg.IP, n.node.IP) {
			fmt.Printf("updated node: %s, %s\n", cfg.Name, cfg.IP.String())
			c.Nodes[i].node = cfg
			c.Nodes[i].lastSeen = time.Now()
			return nil
		}
	}
	fmt.Printf("added node: %s, %s\n", cfg.Name, cfg.IP.String())
	c.Nodes = append(c.Nodes, controlNode{node: cfg, lastSeen: time.Now()})

	return nil
}

// deleteNode will delete a Node from the list of known nodes
func (c *Controller) deleteNode(node Config) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for i, n := range c.Nodes {
		if bytes.Equal(node.IP, n.node.IP) {
			c.Nodes = append(c.Nodes[:i], c.Nodes[i+1:]...)
		}
	}

	return fmt.Errorf("no known node with this ip known, ip: %s", node.IP)
}

// gcNode will remove stale Nodes from the list of known nodes
func (c *Controller) gcNode() {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

start:
	for i := range c.Nodes {
		if c.Nodes[i].lastSeen.Add(10 * time.Second).Before(time.Now()) {
			// it has been more then 10 seconds since we saw this node. remove it now.
			fmt.Printf("remove stale node: %s, %s\n", c.Nodes[i].node.Name, c.Nodes[i].node.IP.String())
			c.Nodes = append(c.Nodes[:i], c.Nodes[i+1:]...)
			goto start
		}
	}
}
