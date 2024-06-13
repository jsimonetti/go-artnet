package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jsimonetti/go-artnet"
	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/types"
)

func main() {
	ip := mustGetLocalAddress().IP

	log := artnet.NewDefaultLogger()
	opt := artnet.BroadcastAddr(*mustGetBroadcastAddress())
	c := artnet.NewController("controller-1", ip, log, opt)
	c.Start(context.Background())

	time.Sleep(time.Second * 1)

	c.LogNodes()
	time.Sleep(time.Second * 1)

	c.LogNodes()
	time.Sleep(time.Second * 1)

	c.RangeIPs(func(ip string) {
		log.With(artnet.Fields{"ip": ip}).Infof("Range over ip")
		c.RangeOutputsOf(ip, func(a types.Address) {
			log.With(artnet.Fields{"address": a, "ip": ip}).Infof("Range over address")
			for i := 0; i < 30; i++ {
				err := c.SendDMX(net.ParseIP(ip), a, get170(0xff, 0x00, 0x00))
				if err != nil {
					log.With(artnet.Fields{"error": err}).Error("Failed to Send DMX Data")
				}
				time.Sleep(time.Millisecond * 30)
			}

			for i := 0; i < 30; i++ {
				err := c.SendDMX(net.ParseIP(ip), a, get170(0x00, 0xff, 0x00))
				if err != nil {
					log.With(artnet.Fields{"error": err}).Error("Failed to Send DMX Data")
				}
				time.Sleep(time.Millisecond * 30)
			}

			for i := 0; i < 30; i++ {
				err := c.SendDMX(net.ParseIP(ip), a, get170(0x00, 0x00, 0xff))
				if err != nil {
					log.With(artnet.Fields{"error": err}).Error("Failed to Send DMX Data")
				}
				time.Sleep(time.Millisecond * 30)
			}
		})
	})

	time.Sleep(time.Second * 60)

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
