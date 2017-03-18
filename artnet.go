package artnet

import (
	"net"

	"github.com/jsimonetti/go-artnet/node"
	"github.com/jsimonetti/go-artnet/packet/code"
)

// NewController returns an instance of an ArtNet Controller
func NewController(name string, ip net.IP) *node.Controller {
	return node.NewController(name, ip)
}

// NewNode returns an instance of an ArtNet Node
func NewNode(name string, style code.StyleCode, ip net.IP) *node.Node {
	return node.NewNode(name, style, ip)
}
