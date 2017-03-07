package packet

import (
	"bytes"
	"encoding"
	"errors"

	"github.com/jsimonetti/go-artnet/packet/code"
)

// Various errors which may occur when attempting to marshal or unmarshal
// an ArtNetPacket to and from its binary form.
var (
	errIncorrectHeaderLength = errors.New("header length incorrect")
	errInvalidPacket         = errors.New("invalid Art-Net packet")
	errInvalidOpCode         = errors.New("invalid OpCode in packet")
)

// ArtNetPacket is the interface used for passing around different kinds of ArtNet packets.
type ArtNetPacket interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	artPacket()
}

// ArtNet is the fixed string "Art-Net" terminated with a zero
var ArtNet = [8]byte{0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00}

var _ ArtNetPacket = &Packet{}

// Packet contains the base header for a ArtNet Packet
type Packet struct {
	// ID is an Array of 8 characters, the final character is a null termination.
	// Value should be []byte{‘A’,‘r’,‘t’,‘-‘,‘N’,‘e’,‘t’,0x00}
	//ID     []byte

	// OpCode defines the class of data following within this UDP packet.
	// Transmitted low byte first.
	OpCode code.OpCode

	// data contains the rest of the data for this packet
	data []byte
}

// MarshalBinary marshals a Packet into a byte slice.
func (p *Packet) MarshalBinary() (b []byte, err error) {
	b = make([]byte, 10)
	copy(b[0:8], ArtNet[:]) // ID is always fixed

	if !p.OpCode.Valid() {
		return nil, errInvalidOpCode
	}
	copy(b[8:10], p.OpCode.Marshal())
	return append(b, p.data...), nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into a Packet.
func (p *Packet) UnmarshalBinary(b []byte) error {
	if len(b) < 10 {
		return errIncorrectHeaderLength
	}
	if !bytes.Equal(b[0:8], ArtNet[:]) {
		return errInvalidPacket
	}

	p.OpCode = code.OpCode.Unmarshal(p.OpCode, b[8:10])

	if !p.OpCode.Valid() {
		return errInvalidOpCode
	}

	p.data = append(p.data, b[10:]...)

	return nil
}

// artPacket is an empty method to sattisfy the ArtNetPacket interface.
func (p *Packet) artPacket() {}
