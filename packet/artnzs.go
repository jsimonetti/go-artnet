package packet

import (
	"github.com/jsimonetti/go-artnet/packet/code"
)

var _ ArtNetPacket = &ArtNzsPacket{}

// ArtNzsPacket contains an ArtNzs Packet.
//
// ArtNzs is the data packet used to transfer DMX512 data with non-zero start codes
// (except RDM). The format is identical for Node to Controller, Node to Node and
// Controller to Node
//
// Packet Strategy:
//  Controller -  Receive:            Application Specific
//                Unicast Transmit:   Yes
//                Broadcast Transmit: No
//  Node -        Receive:            Application Specific
//                Unicast Transmit:   Yes
//                Broadcast Transmit: No
//  MediaServer - Receive:            Application Specific
//                Unicast Transmit:   Yes
//                Broadcast Transmit: No
type ArtNzsPacket struct {
	// Inherit the Header header
	Header

	// Sequence number is used to ensure that ArtNzs packets are used in the correct order.
	// When Art-Net is carried over a medium such as the Internet, it is possible that ArtNzs packets
	// will reach the receiver out of order.
	// This field is incremented in the range 0x01 to 0xff to allow the receiving node to resequence
	// packets. The Sequence field is set to 0x00 to disable this feature
	Sequence uint8

	// StartCode is the DMX512 start code of this packet. Must not be Zero or RDM
	StartCode uint8

	// SubUni is the low byte of the 15 bit Port-Address to which this packet is destined
	SubUni uint8

	// Net is the top 7 bits of the 15 bit Port-Address to which this packet is destined
	Net uint8

	// Length indicates the length of the data. This value should be a number in the
	// range 1 – 512. It represents the number of DMX512 channels encoded in packet.
	Length uint16

	// Data is a variable length string of DMX512 lighting data
	Data [512]byte
}

// NewArtNzsPacket returns an ArtNetPacket with the correct OpCode
func NewArtNzsPacket() *ArtNzsPacket {
	return &ArtNzsPacket{}
}

// MarshalBinary marshals an ArtNzsPacket into a byte slice.
func (p *ArtNzsPacket) MarshalBinary() ([]byte, error) {
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtNzsPacket.
func (p *ArtNzsPacket) UnmarshalBinary(b []byte) error {
	return unmarshalPacket(p, b)
}

// validate is used to validate the Packet.
func (p *ArtNzsPacket) validate() error {
	if err := p.Header.validate(); err != nil {
		return err
	}
	if p.OpCode != code.OpNzs {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtNzsPacket) finish() {
	p.Header.finish()
}
