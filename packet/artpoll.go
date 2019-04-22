package packet

import (
	"github.com/jsimonetti/go-artnet/packet/code"
)

var _ ArtNetPacket = &ArtPollPacket{}

// ArtPollPacket contains an ArtPoll Packet.
//
// The ArtPoll packet is used to discover the presence of other Controllers, Nodes and Media Servers.
// The ArtPoll packet is only sent by a Controller. Both Controllers and Nodes respond to the packet.
//
// A Controller broadcasts an ArtPoll packet to IP address 2.255.255.255 (sub-net mask 255.0.0.0) at
// UDP port 0x1936, this is the Directed Broadcast address.
//
// The Controller may assume a maximum timeout of 3 seconds between sending ArtPoll and receiving all
// ArtPollReply packets. If the Controller does not receive a response in this time it should consider
// the Node to have disconnected.
//
// The Controller that broadcasts an ArtPoll should also reply to its own message (to
// Directed Broadcast address) with an ArtPollReply. This ensures that any other Controllers listening
// to the network will detect all devices without the need for all Controllers connected to the
// network to send ArtPoll packets. It is a requirement of Art-Net that all controllers broadcast
// an ArtPoll every 2.5 to 3 seconds. This ensures that any network devices can easily detect a
// disconnect.
//
// Packet Strategy:
//  Controller -  Receive:            Send ArtPollReply
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Controller broadcasts this packet to poll all Controllers and
//                                   Nodes on the network
//  Node -        Receive:            Send ArtPollReply
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Not Allowed
//  MediaServer - Receive:            Send ArtPollReply
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Not Allowed
type ArtPollPacket struct {
	// Inherit the Header header
	Header

	// TalkToMe defines the behavior of the Node
	TalkToMe code.TalkToMe

	// Priority contains the lowest priority of diagnostics message that should be sent
	Priority code.PriorityCode
}

// NewArtPollPacket returns an ArtNetPacket with the correct OpCode
func NewArtPollPacket() *ArtPollPacket {
	return &ArtPollPacket{}
}

// MarshalBinary marshals an ArtPollPacket into a byte slice.
func (p *ArtPollPacket) MarshalBinary() ([]byte, error) {
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollPacket.
func (p *ArtPollPacket) UnmarshalBinary(b []byte) error {
	return unmarshalPacket(p, b)
}

// validate is used to validate the Packet.
func (p *ArtPollPacket) validate() error {
	if err := p.Header.validate(); err != nil {
		return err
	}
	if p.OpCode != code.OpPoll {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtPollPacket) finish() {
	p.OpCode = code.OpPoll
	p.Header.finish()
}
