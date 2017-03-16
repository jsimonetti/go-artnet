package main

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/packet"
)

func main() {

	dst := fmt.Sprintf("%s:%d", "255.255.255.255", packet.ArtNetPort)
	broadcastAddr, _ := net.ResolveUDPAddr("udp", dst)
	src := fmt.Sprintf("%s:%d", "2.12.12.12", packet.ArtNetPort)
	localAddr, _ := net.ResolveUDPAddr("udp", src)

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Printf("error opening udp: %s\n", err)
		return
	}

	// set channels 1, 2, 3 and 4 to FL

	p := &packet.ArtDMXPacket{
		Sequence: 1,
		SubUni:   0,
		Net:      0,
		Data:     [512]byte{0xff, 0xff, 0xff, 0xff},
	}

	b, err := p.MarshalBinary()

	n, err := conn.WriteTo(b, broadcastAddr)
	if err != nil {
		fmt.Printf("error writing packet: %s\n", err)
		return
	}
	fmt.Printf("packet sent, wrote %d bytes\n", n)
}
