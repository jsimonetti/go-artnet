package main

import (
	"fmt"
	"net"
	"time"

	"github.com/jsimonetti/go-artnet/packet/code"

	"github.com/jsimonetti/go-artnet"
)

func main() {

	artsubnet := "2.0.0.0/8"
	_, cidrnet, _ := net.ParseCIDR(artsubnet)

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("error getting ips: %s\n", err)
	}

	var ip net.IP

	for _, addr := range addrs {
		ip = addr.(*net.IPNet).IP
		if cidrnet.Contains(ip) {
			break
		}
	}

	log := artnet.NewDefaultLogger()
	n := artnet.NewNode("node-1", code.StNode, ip, log)
	n.Start()

	for {
		time.Sleep(time.Second)
	}
}
