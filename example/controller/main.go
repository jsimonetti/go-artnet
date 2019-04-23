package main

import (
	"fmt"
	"net"
	"time"

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
	c := artnet.NewController("controller-1", ip, log)
	c.Start()

	go func() {
		time.Sleep(10 * time.Second)
		c.SendDMXToAddress([512]byte{0x00, 0xff, 0x00, 0xff, 0x00}, artnet.Address{Net: 0, SubUni: 0})
		time.Sleep(5 * time.Second)
		c.SendDMXToAddress([512]byte{0xff, 0x00, 0x00, 0xff, 0x00}, artnet.Address{Net: 0, SubUni: 0})
		time.Sleep(5 * time.Second)
		c.SendDMXToAddress([512]byte{0x00, 0x00, 0xff, 0xff, 0x00}, artnet.Address{Net: 0, SubUni: 0})
		time.Sleep(5 * time.Second)
		c.SendDMXToAddress([512]byte{}, artnet.Address{Net: 0, SubUni: 0})
		time.Sleep(5 * time.Second)
	}()

	for {
		time.Sleep(time.Second)
	}
}
