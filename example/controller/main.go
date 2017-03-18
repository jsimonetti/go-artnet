package main

import (
	"net"

	"github.com/jsimonetti/go-artnet"
)

func main() {
	c := artnet.NewController("controller-1", net.ParseIP("2.12.12.12"))
	if err := c.Start(); err != nil {
		panic(err)
	}
}
