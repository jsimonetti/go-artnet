package packet

import (
	"bytes"
	"encoding/binary"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

var _ ArtNetPacket = &ArtPollPacket{}

// ArtPollPacket contains an ArtPoll Packet.
type ArtPollPacket struct {
	// Inherit the Packet header
	Packet

	// this packet type contains a version
	version [2]byte

	// TalkToMe defines the behavior of the Node
	TalkToMe code.TalkToMe

	// Priority contains the lowest priority of diagnostics message that should be sent
	Priority code.PriorityCode
}

// NewArtPollPacket returns an ArtNetPacket with the correct OpCode
func NewArtPollPacket() *ArtPollPacket {
	return &ArtPollPacket{
		Packet: Packet{
			OpCode: code.OpPoll,
			id:     ArtNet,
		},
		version: version.Bytes(),
	}
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *ArtPollPacket) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), p.validate()
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
//TODO
func (p *ArtPollPacket) UnmarshalBinary(b []byte) error {
	return p.validate()
}

// artPacket is an empty method to sattisfy the ArtNetPacket interface.
func (p *ArtPollPacket) validate() error {
	if p.OpCode != code.OpPoll {
		return errInvalidOpCode
	}
	return nil
}
