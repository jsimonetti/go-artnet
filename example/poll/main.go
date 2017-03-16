package main

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/packet"
)

func main() {

	dst := fmt.Sprintf("%s:%d", "255.255.255.255", packet.ArtNetPort)
	broadcastAddr, _ := net.ResolveUDPAddr("udp", dst)
	//src := fmt.Sprintf("%s:%d", "2.12.12.12", packet.ArtNetPort)
	localAddr, _ := net.ResolveUDPAddr("udp", "2.12.12.12:6454")

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Printf("error opening udp: %s\n", err)
		return
	}

	p := &packet.ArtPollPacket{}
	b, err := p.MarshalBinary()
	if err != nil {
		fmt.Printf("error marshalling packet: %s\n", err)
		return
	}

	n, err := conn.WriteTo(b, broadcastAddr)
	if err != nil {
		fmt.Printf("error writing packet: %s\n", err)
		return
	}
	fmt.Printf("packet sent, wrote %d bytes\n", n)

	buf := make([]byte, 4096)
	n, a, err := conn.ReadFrom(buf)
	n, a, err = conn.ReadFrom(buf)
	if err != nil {
		fmt.Printf("error reading packet: %s\n", err)
		return
	}
	fmt.Printf("packet read from %v, read %d bytes\n", a, n)

	r := &packet.ArtPollReplyPacket{}
	err = r.UnmarshalBinary(buf)
	if err != nil {
		fmt.Printf("error unmarshalling packet: %s\n", err)
		return
	}
	fmt.Printf("got reply: %v", r)
}
