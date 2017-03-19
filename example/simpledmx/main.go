package main

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/packet"
)

func main() {

	dst := fmt.Sprintf("%s:%d", "2.231.20.36", packet.ArtNetPort)
	node, _ := net.ResolveUDPAddr("udp", dst)
	src := fmt.Sprintf("%s:%d", "2.12.12.12", packet.ArtNetPort)
	localAddr, _ := net.ResolveUDPAddr("udp", src)

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Printf("error opening udp: %s\n", err)
		return
	}

	// set channels 1 and 4 to FL, 2, 3 and 5 to FD
	// on my colorBeam this sets output 1 to fullbright red with zero strobing

	p := &packet.ArtDMXPacket{
		Sequence: 1,
		SubUni:   0,
		Net:      0,
		Data:     [512]byte{0xff, 0x00, 0x00, 0xff, 0x00},
	}

	b, err := p.MarshalBinary()

	n, err := conn.WriteTo(b, node)
	if err != nil {
		fmt.Printf("error writing packet: %s\n", err)
		return
	}
	fmt.Printf("packet sent, wrote %d bytes\n", n)
}
