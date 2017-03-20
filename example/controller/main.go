package main

import (
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"

	artnet "github.com/jsimonetti/go-artnet"
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

	c := artnet.NewController("controller-1", ip)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		c.Start()
		wg.Done()
	}()
	time.Sleep(10 * time.Second)
	c.SendDMXToAddress([512]byte{0x00, 0xff, 0x00, 0xff, 0x00}, artnet.Address{Net: 0, SubUni: 0})
	time.Sleep(5 * time.Second)
	c.SendDMXToAddress([512]byte{0xff, 0x00, 0x00, 0xff, 0x00}, artnet.Address{Net: 0, SubUni: 0})
	time.Sleep(5 * time.Second)
	c.SendDMXToAddress([512]byte{0x00, 0x00, 0xff, 0xff, 0x00}, artnet.Address{Net: 0, SubUni: 0})
	time.Sleep(5 * time.Second)
	c.SendDMXToAddress([512]byte{}, artnet.Address{Net: 0, SubUni: 0})
	time.Sleep(5 * time.Second)
	c.Stop()
	wg.Wait()
	fmt.Printf("num: %d", runtime.NumGoroutine())

}
