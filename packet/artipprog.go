package packet

import (
	"bytes"
	"encoding/binary"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

var _ ArtNetPacket = &ArtIPProgPacket{}

// ArtIPProgPacket contains an ArtIPProg Packet.
//
// The ArtIpProg packet allows the IP settings of a Node to be reprogrammed. The ArtIpProg
// packet is sent by a Controller to the private address of a Node. If the Node supports
// remote programming of IP address, it will respond with an ArtIpProgReply packet.
// In all scenarios, the ArtIpProgReply is sent to the private address of the sender
//
// Packet Strategy:
//  Controller -  Receive:            No Action
//                Unicast Transmit:   Controller transmits to a specific node IP address
//                Broadcast Transmit: Not Allowed
//  Node -        Receive:            Reply with ArtIpProgReply
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Not Allowed
//  MediaServer - Receive:            Reply with ArtIpProgReply
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Not Allowed
type ArtIPProgPacket struct {
	// Inherit the Header header
	Header

	// this packet type contains a version
	version [2]byte

	// filler to pad to same length as ArtPoll
	filler [2]byte

	// Command defines the how this packet is processed. If all bits are clear, this
	// is an enquiry only
	Command uint8

	// filler2 to pad to word allignment
	filler2 byte

	// ProgIP is IP Address to be programmed into Node if enabled by Command Field
	ProgIP [4]byte

	// ProgSubNet is Subnet mask to be programmed into Node if enabled by Command Field
	ProgSubNet [4]byte

	//ProgPort is deprecated
	ProgPort [2]byte

	//spare bytes, transmit as zero, receivers don’t test.
	spare [8]byte
}

// NewArtIPProgPacket returns an ArtNetPacket with the correct OpCode
func NewArtIPProgPacket() *ArtIPProgPacket {
	return &ArtIPProgPacket{
		Header: Header{
			OpCode: code.OpPoll,
			id:     ArtNet,
		},
		version: version.Bytes(),
	}
}

// MarshalBinary marshals an ArtIPProgPacket into a byte slice.
func (p *ArtIPProgPacket) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), p.validate()
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtIPProgPacket.
//TODO
func (p *ArtIPProgPacket) UnmarshalBinary(b []byte) error {
	return p.validate()
}

// artPacket is an empty method to sattisfy the ArtNetPacket interface.
func (p *ArtIPProgPacket) validate() error {
	if p.OpCode != code.OpIPProg {
		return errInvalidOpCode
	}
	return nil
}
