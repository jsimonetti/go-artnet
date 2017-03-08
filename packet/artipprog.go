package packet

import (
	"github.com/jsimonetti/go-artnet/packet/code"
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

	// Filler1 to pad to same length as ArtPoll
	_ [2]byte

	// Command defines the how this packet is processed. If all bits are clear, this
	// is an enquiry only
	Command uint8

	// Filler2 to pad to word allignment
	_ byte

	// ProgIP is IP Address to be programmed into Node if enabled by Command Field
	ProgIP [4]byte

	// ProgSubNet is Subnet mask to be programmed into Node if enabled by Command Field
	ProgSubNet [4]byte

	// ProgPort is deprecated
	ProgPort [2]byte

	// Spare bytes, transmit as zero, receivers don’t test.
	_ [8]byte
}

// NewArtIPProgPacket returns an ArtNetPacket with the correct OpCode
func NewArtIPProgPacket() *ArtIPProgPacket {
	return &ArtIPProgPacket{}
}

// MarshalBinary marshals an ArtIPProgPacket into a byte slice.
func (p *ArtIPProgPacket) MarshalBinary() ([]byte, error) {
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtIPProgPacket.
func (p *ArtIPProgPacket) UnmarshalBinary(b []byte) error {
	return unmarshalPacket(p, b)
}

// validate is used to validate the Packet.
func (p *ArtIPProgPacket) validate() error {
	if err := p.Header.validate(); err != nil {
		return err
	}
	if p.OpCode != code.OpIPProg {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtIPProgPacket) finish() {
	p.Header.finish()
}
