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
	errInvalidStyle          = errors.New("invalid StyleCode in packet")
)

// ArtNetPacket is the interface used for passing around different kinds of ArtNet packets.
type ArtNetPacket interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	validate() error
}

// ArtNet is the fixed string "Art-Net" terminated with a zero
var ArtNet = [8]byte{0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00}

// Header contains the base header for an ArtNet Packet
type Header struct {
	// ID is an Array of 8 characters, the final character is a null termination.
	// Value should be []byte{‘A’,‘r’,‘t’,‘-‘,‘N’,‘e’,‘t’,0x00}
	id [8]byte

	// OpCode defines the class of data following within this UDP packet.
	// Transmitted low byte first.
	OpCode code.OpCode
}

// UnmarshalBinary unmarshals the contents of a byte slice into a Packet.
func (p *Header) unmarshal(b []byte) error {
	if len(b) < 10 {
		return errIncorrectHeaderLength
	}
	if !bytes.Equal(b[0:8], ArtNet[:]) {
		return errInvalidPacket
	}
	p.OpCode = code.OpCode(uint16(b[8] + b[9]))
	return p.validate()
}

func (p *Header) validate() error {
	if !code.ValidOp(p.OpCode) {
		return errInvalidOpCode
	}
	return nil
}
