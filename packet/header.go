package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/jsimonetti/go-artnet/artnettypes"
	"github.com/jsimonetti/go-artnet/packet/code"
)

// Various errors which may occur when attempting to marshal or unmarshal
// an ArtNetPacket to and from its binary form.
var (
	errIncorrectHeaderLength = errors.New("header length incorrect")
	errInvalidPacket         = errors.New("invalid Art-Net packet")
	errInvalidPacketBoundary = errors.New("invalid Art-Net packet, not aligned on 16 bit boundary")
	errInvalidPacketMin      = errors.New("invalid Art-Net packet, less than min")
	errInvalidPacketMax      = errors.New("invalid Art-Net packet, greater than max")
	errInvalidOpCode         = errors.New("invalid OpCode in packet")
	errInvalidStyleCode      = errors.New("invalid StyleCode in packet")
	errNotImplementedOpCode  = errors.New("not implemented OpCode in packet")
)

// ArtNetPort is the fixed ArtNet port 6454.
const ArtNetPort = 6454

const version = artnettypes.CurrentVersion

var artNet = artnettypes.ArtNet

// Header contains the base header for an ArtNet Packet
type Header struct {
	// ID is an Array of 8 characters, the final character is a null termination.
	// Value should be []byte{‘A’,‘r’,‘t’,‘-‘,‘N’,‘e’,‘t’,0x00}
	artnettypes.ID

	// OpCode defines the class of data following within this UDP packet.
	// Transmitted low byte first, GetOpCode() addresses this.
	code.OpCode

	// Version of this packet
	// Transmitted low byte first, GetVersion() addresses this.
	artnettypes.Version
}

func NewHeader(opcode code.OpCode) Header {
	return Header{
		ID:      artNet,
		OpCode:  code.OpCode(swapUint16(uint16(opcode))),
		Version: version,
	}
}

func (p *Header) validate(expectedOpCode code.OpCode) error {
	if p.ID != artnettypes.ArtNet {
		return errInvalidPacket
	}

	if p.GetOpCode() != expectedOpCode {
		return errInvalidOpCode
	}

	if p.GetVersion() < version {
		return fmt.Errorf("incompatible version. want: %d, got: %d", version, p.GetVersion())
	}

	return nil
}

// unmarshal the contents of a byte slice into a Header.
func (h *Header) unmarshal(b []byte) error {
	buf := bytes.NewReader(b)
	return binary.Read(buf, binary.BigEndian, h)
}

// HeaderWithoutVersion contains the base header for an ArtNet Packet
// without the version
type HeaderWithoutVersion struct {
	// ID is an Array of 8 characters, the final character is a null termination.
	// Value should be []byte{‘A’,‘r’,‘t’,‘-‘,‘N’,‘e’,‘t’,0x00}
	artnettypes.ID

	// OpCode defines the class of data following within this UDP packet.
	// Transmitted low byte first.
	code.OpCode
}

func NewHeaderWithoutVersion(opcode code.OpCode) HeaderWithoutVersion {
	return HeaderWithoutVersion{
		ID:     artnettypes.ArtNet,
		OpCode: code.OpCode(swapUint16(uint16(opcode))),
	}
}

func (p *HeaderWithoutVersion) validate(expectedOpCode code.OpCode) error {
	if p.ID != artnettypes.ArtNet {
		return errInvalidPacket
	}

	if p.GetOpCode() != expectedOpCode {
		return errInvalidOpCode
	}

	return nil
}
