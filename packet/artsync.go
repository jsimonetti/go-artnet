package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

var _ ArtNetPacket = &ArtSyncPacket{}

// ArtSyncPacket contains an ArtSync Packet.
//
// The ArtSync packet can be used to force nodes to synchronously output ArtDmx packets
// to their outputs. This is useful in video and media-wall applications.
// A controller that wishes to implement synchronous transmission will unicast multiple
// universes of ArtDmx and then broadcast an ArtSync to synchronously transfer all the
// ArtDmx packets to the nodes’ outputs at the same time.
//
// Packet Strategy:
//  Controller -  Receive:            No Action
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Controller broadcasts this packet to synchronously
//                                    transfer previous ArtDmx packets to Node’s output
//  Node -        Receive:            Transfer previous ArtDmx packets to output
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Not Allowed
//  MediaServer - Receive:            Transfer previous ArtDmx packets to output
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Not Allowed
type ArtSyncPacket struct {
	// Inherit the Header header
	Header

	// this packet type contains a version
	version [2]byte

	// aux are auxiliary bytes transmitted as zero
	aux [2]byte
}

// NewArtSyncPacket returns an ArtNetPacket with the correct OpCode
func NewArtSyncPacket() *ArtSyncPacket {
	return &ArtSyncPacket{}
}

// MarshalBinary marshals an ArtSyncPacket into a byte slice.
func (p *ArtSyncPacket) MarshalBinary() ([]byte, error) {
	p.finish()
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtSyncPacket.
func (p *ArtSyncPacket) UnmarshalBinary(b []byte) error {
	if err := p.Header.unmarshal(b[:10]); err != nil {
		return err
	}
	if len(b) != 14 {
		return fmt.Errorf("invalid packet length received. want: 14, got: %d", len(b))
	}
	p.version = [2]byte{b[11], b[10]}
	p.aux = [2]byte{b[12], b[13]}

	return p.validate()
}

// validate is used to validate the Packet.
func (p *ArtSyncPacket) validate() error {
	if p.OpCode != code.OpSync {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtSyncPacket) finish() {
	p.OpCode = code.OpSync
	p.version = version.Bytes()
	p.id = ArtNet
}
