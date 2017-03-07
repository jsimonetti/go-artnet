package packet

import (
	"bytes"
	"encoding/binary"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
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

	// this packet type contains a version
	version [2]byte

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
	Data string
}

// NewArtNzsPacket returns an ArtNetPacket with the correct OpCode
func NewArtNzsPacket() *ArtNzsPacket {
	return &ArtNzsPacket{
		Header: Header{
			OpCode: code.OpNzs,
			id:     ArtNet,
		},
		version: version.Bytes(),
	}
}

// MarshalBinary marshals an ArtNzsPacket into a byte slice.
func (p *ArtNzsPacket) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), p.validate()
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtNzsPacket.
//TODO
func (p *ArtNzsPacket) UnmarshalBinary(b []byte) error {
	return p.validate()
}

// artPacket is an empty method to sattisfy the ArtNetPacket interface.
func (p *ArtNzsPacket) validate() error {
	if p.OpCode != code.OpNzs {
		return errInvalidOpCode
	}
	return nil
}
