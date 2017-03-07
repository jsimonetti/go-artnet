package packet

import (
	"bytes"
	"encoding/binary"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

var _ ArtNetPacket = &ArtTimeCodePacket{}

// ArtTimeCodePacket contains an ArtTimeCode Packet.
//
// ArtTimeCode allows time code to be transported over the network. The data format is
// compatible with both longitudinal time code and MIDI time code. The four key types of
// Film, EBU, Drop Frame and SMPTE are also encoded.
// Use of the packet is application specific but in general a single controller will
// broadcast the packet to the network.
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
type ArtTimeCodePacket struct {
	// Inherit the Header header
	Header

	// filler1
	filler1 [2]byte

	// Frames time. 0 – 29 depending on mode
	Frames uint8

	// Seconds 0 - 59
	Seconds uint8

	// Minutes 0 - 59
	Minutes uint8

	//Hours 0 - 23
	Hours uint8

	// Type of source, 0 = Film (24fps), 1 = EBU (25fps), 2 = DF (29.97fps), 3 = SMPTE (30fps)
	Type uint8
}

// NewArtTimeCodePacket returns an ArtNetPacket with the correct OpCode
func NewArtTimeCodePacket() *ArtTimeCodePacket {
	return &ArtTimeCodePacket{}
}

// MarshalBinary marshals an ArtTimeCodePacket into a byte slice.
func (p *ArtTimeCodePacket) MarshalBinary() ([]byte, error) {
	p.finish()
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtTimeCodePacket.
//TODO
func (p *ArtTimeCodePacket) UnmarshalBinary(b []byte) error {
	return p.validate()
}

// validate is used to validate the Packet.
func (p *ArtTimeCodePacket) validate() error {
	if p.OpCode != code.OpTimeCode {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtTimeCodePacket) finish() {
	p.OpCode = code.OpTimeCode
	p.version = version.Bytes()
	p.id = ArtNet
}
