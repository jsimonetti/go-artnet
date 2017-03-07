package packet

import (
	"bytes"
	"encoding/binary"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

var _ ArtNetPacket = &ArtIPProgReplyPacket{}

// ArtIPProgReplyPacket contains an ArtIPProgReply Packet.
//
// The ArtIpProgReply packet is issued by a Node in response to an ArtIpProg packet.
// Nodes that do not support remote programming of IP address do not reply to ArtIpProg
// packets. In all scenarios, the ArtIpProgReply is sent to the private address of the
// sender.
//
// Packet Strategy:
//  Controller -  Receive:            No Action
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Not Allowed
//  Node -        Receive:            No Action
//                Unicast Transmit:   Transmits to specific Controller IP address
//                Broadcast Transmit: Not Allowed
//  MediaServer - Receive:            No Action
//                Unicast Transmit:   Transmits to specific Controller IP address
//                Broadcast Transmit: Not Allowed
type ArtIPProgReplyPacket struct {
	// Inherit the Header header
	Header

	// filler to pad to same length as ArtPoll and ArtIpProg
	filler [4]byte

	// ProgIP is IP Address to be programmed into Node if enabled by Command Field
	ProgIP [4]byte

	// ProgSubNet is Subnet mask to be programmed into Node if enabled by Command Field
	ProgSubNet [4]byte

	// ProgPort is deprecated
	ProgPort [2]byte

	// Status defines if DHCP is enabled or not
	Status uint8

	// spare bytes, transmit as zero, receivers don’t test.
	spare [7]byte
}

// NewArtIPProgReplyPacket returns an ArtNetPacket with the correct OpCode
func NewArtIPProgReplyPacket() *ArtIPProgReplyPacket {
	return &ArtIPProgReplyPacket{}
}

// MarshalBinary marshals an ArtIPProgReplyPacket into a byte slice.
func (p *ArtIPProgReplyPacket) MarshalBinary() ([]byte, error) {
	p.finish()
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtIPProgReplyPacket.
//TODO
func (p *ArtIPProgReplyPacket) UnmarshalBinary(b []byte) error {
	return p.validate()
}

// validate is used to validate the Packet.
func (p *ArtIPProgReplyPacket) validate() error {
	if p.OpCode != code.OpIPProgReply {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtIPProgReplyPacket) finish() {
	p.OpCode = code.OpIPProgReply
	p.id = ArtNet
	p.version = version.Bytes()
}
