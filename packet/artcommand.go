package packet

import (
	"bytes"
	"encoding/binary"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

var _ ArtNetPacket = &ArtCommandPacket{}

// ArtCommandPacket contains an ArtCommand Packet.
//
// The ArtCommand packet is used to send property set style commands. The packet can be
// unicast or broadcast, the decision being application specific.
//
// Packet Strategy:
//  Controller -  Receive:            Application Specific
//                Unicast Transmit:   Application Specific
//                Broadcast Transmit: Application Specific
//  Node -        Receive:            Application Specific
//                Unicast Transmit:   Application Specific
//                Broadcast Transmit: Application Specific
//  MediaServer - Receive:            Application Specific
//                Unicast Transmit:   Application Specific
//                Broadcast Transmit: Application Specific
type ArtCommandPacket struct {
	// Inherit the Header header
	Header

	// this packet type contains a version
	version [2]byte

	// ESTAmanufacturer contains a code used to represent equipment manufacturer.
	// They are assigned by ESTA. This field can be interpreted as two ASCII bytes
	// representing the manufacturer initials.
	ESTAmanufacturer [2]byte

	// Length indicates the length of the data
	Length uint16

	// Data is an ASCII string, null terminated. Max length is 512 bytes including the null terminator
	Data string
}

// NewArtCommandPacket returns an ArtNetPacket with the correct OpCode
func NewArtCommandPacket() *ArtCommandPacket {
	return &ArtCommandPacket{
		Header: Header{
			OpCode: code.OpCommand,
			id:     ArtNet,
		},
		version: version.Bytes(),
	}
}

// MarshalBinary marshals an ArtCommandPacket into a byte slice.
func (p *ArtCommandPacket) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), p.validate()
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtCommandPacket.
//TODO
func (p *ArtCommandPacket) UnmarshalBinary(b []byte) error {
	return p.validate()
}

// artPacket is an empty method to sattisfy the ArtNetPacket interface.
func (p *ArtCommandPacket) validate() error {
	if p.OpCode != code.OpCommand {
		return errInvalidOpCode
	}
	return nil
}
