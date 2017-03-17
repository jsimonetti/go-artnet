package node

import (
	"bytes"
	"fmt"
	"sync"
)

// Controller holds the information for a controller
type Controller struct {
	// Node is the controller itself
	Node

	// Nodes is a slice of nodes that are seen by this controller
	Nodes    []NodeConfig
	nodeLock sync.Mutex
}

func (c *Controller) Start() {
}

func (c *Controller) pollLoop() {
}

// addNode will add a Node to the list of known nodes
func (c *Controller) addNode(node NodeConfig) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for _, n := range c.Nodes {
		if bytes.Equal(node.IP, n.IP) {
			return fmt.Errorf("allready a Node with this ip known, ip: %s", node.IP)
		}
	}
	c.Nodes = append(c.Nodes, node)

	return nil
}

// deleteNode will delete a Node from the list of known nodes
func (c *Controller) deleteNode(node NodeConfig) error {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for i, n := range c.Nodes {
		if bytes.Equal(n.IP, node.IP) {
			c.Nodes = append(c.Nodes[:i], c.Nodes[i+1:]...)
		}
	}

	return fmt.Errorf("no known node with this ip known, ip: %s", node.IP)
}
