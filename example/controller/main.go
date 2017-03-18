package main

import (
	"net"

	"github.com/jsimonetti/go-artnet/node"
	"github.com/jsimonetti/go-artnet/packet/code"
)

func main() {
	c := &node.Controller{}
	c.Node = node.New("controller-1", code.StController, net.ParseIP("2.12.12.12"))
	if err := c.Start(); err != nil {
		panic(err)
	}
}
