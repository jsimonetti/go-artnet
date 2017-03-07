package packet

import (
	"github.com/jsimonetti/artnet/packet/code"
	"github.com/jsimonetti/artnet/version"
)

var _ ArtNetPacket = &ArtPollPacket{}

// ArtPollPacket contains an ArtPoll Packet.
type ArtPollPacket struct {
	// Inherit the Packet header
	Packet

	// TalkToMe defines the behavior of the Node
	TalkToMe TalkToMe

	// Priority contains the lowest priority of diagnostics message that should be sent
	Priority code.PriorityCode
}

// NewArtPollPacket returns an ArtNetPacket with the correct OpCode
func NewArtPollPacket() *ArtPollPacket {
	return &ArtPollPacket{}
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *ArtPollPacket) MarshalBinary() ([]byte, error) {
	p.Packet.OpCode = code.OpPoll
	p.Packet.data = append(p.Packet.data, version.Bytes()...)
	p.Packet.data = append(p.Packet.data, uint8(p.TalkToMe), uint8(p.Priority))

	return p.Packet.MarshalBinary()
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
//TODO
func (p *ArtPollPacket) UnmarshalBinary(b []byte) error {
	return nil
}
