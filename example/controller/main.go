package main

import (
	"net"
	"sync"
	"time"

	"github.com/jsimonetti/go-artnet"
)

func main() {
	c := artnet.NewController("controller-1", net.ParseIP("2.12.12.12"))
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		c.Start()
		wg.Done()
	}()
	time.Sleep(10 * time.Second)
	c.SendDMXToAddress([512]byte{0xff, 0x00, 0x00, 0xff, 0x00}, artnet.Address{Net: 0, SubUni: 0})
	time.Sleep(5 * time.Second)
	c.SendDMXToAddress([512]byte{0x00, 0x00, 0xff, 0xff, 0x00}, artnet.Address{Net: 0, SubUni: 0})
	time.Sleep(5 * time.Second)
	c.SendDMXToAddress([512]byte{}, artnet.Address{Net: 0, SubUni: 0})
	time.Sleep(1 * time.Second)
	c.Stop()
	wg.Wait()
}
