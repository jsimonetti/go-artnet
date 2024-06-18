package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jmbarzee/show/common/color"
	"github.com/jsimonetti/go-artnet"
	"github.com/jsimonetti/go-artnet/artnettypes"
	"github.com/jsimonetti/go-artnet/packet"
)

func main() {
	ip := mustGetLocalAddress().IP

	log := artnet.NewDefaultLogger()
	opt := artnet.BroadcastAddr(*mustGetBroadcastAddress())
	c := artnet.NewController("controller-1", ip, log, opt)
	c.Start(context.Background())

	time.Sleep(time.Second * 1)

	// c.LogNodes()
	// time.Sleep(time.Second * 1)

	// c.LogNodes()
	// time.Sleep(time.Second * 1)

	col := color.Blue
	col.SetLightness(0.1)
	for {
		col.ShiftHue(.0004)
		wg := sync.WaitGroup{}
		c.RangeAll(func(ip string, a artnettypes.Address) {
			wg.Add(1)
			r, g, b, _ := col.RGB().ToBytesRGBW()
			// log.With(artnet.Fields{"address": a, "ip": ip, "r": r, "g": g, "b": b}).Infof("Range over address")
			err := c.SendDMX(net.ParseIP(ip), a, get170X(r, g, b))
			if err != nil {
				log.With(artnet.Fields{"error": err}).Error("Failed to Send DMX Data")
			}
			// time.Sleep(time.Millisecond * 30)
			wg.Done()
		})
		time.Sleep(time.Millisecond * 5)
		wg.Wait()
	}

}

func get512(b byte) [512]byte {
	bs := [512]byte{}
	for i := 0; i < 510; i++ {
		bs[i] = b
	}
	return bs
}

func get170(r, g, b byte) [512]byte {
	bs := [512]byte{}
	for i := 0; i < 510; i += 3 {
		bs[i] = r
		bs[i+1] = g
		bs[i+2] = b
	}
	return bs
}

func get170X(r, g, b byte) [512]byte {
	bs := [512]byte{}
	for i := 0; i < 510; i += 3 {
		bs[i] = byte(uint8(r) * uint8(i/37))
		bs[i+1] = byte(uint8(g) * uint8(i/17))
		bs[i+2] = byte(uint8(b) * uint8(i/11))
	}
	return bs
}
func mustGetLocalAddress() *net.UDPAddr {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localIP := conn.LocalAddr().(*net.UDPAddr).IP
	src := fmt.Sprintf("%s:%d", localIP.String(), packet.ArtNetPort)
	fmt.Println(src)

	addr, err := net.ResolveUDPAddr("udp", src)
	if err != nil {
		panic(err)
	}
	return addr
}

func mustGetBroadcastAddress() *net.UDPAddr {
	dst := fmt.Sprintf("%s:%d", "255.255.255.255", packet.ArtNetPort)
	addr, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		panic(err)
	}
	return addr
}
