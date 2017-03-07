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
// The Data field contains the command text. The text is ASCII encoded and is null terminated
// and is case insensitive. It is legal, although inefficient, to set the Data array size to
// the maximum of 512 and null pad unused entries.
// The command text may contain multiple commands and adheres to the following syntax:
//
//   Command=Data&
//
// The ampersand is a break between commands. Also note that the text is capitalised for
// readability; it is case insensitive. Thus far, two commands are defined by Art-Net. It is
// anticipated that additional commands will be added as other manufacturers register commands
// which have industry wide relevance. These commands shall be transmitted with EstaMan = 0xFFFF.
//
// SwoutText - This command is used to re-programme the label associated with the
//             ArtPollReply->Swout fields. Syntax: "SwoutText=Playback&"
// SwinText  - This command is used to re-programme the label associated with the
//             ArtPollReply->Swin fields. Syntax: "SwinText=Record&"
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

	// estamanufacturer contains a code used to represent equipment manufacturer.
	estamanufacturer [2]byte

	// Length indicates the length of the data
	Length uint16

	// Data is an ASCII string, null terminated. Max length is 512 bytes including the null terminator
	Data string
}

// NewArtCommandPacket returns an ArtNetPacket with the correct OpCode
func NewArtCommandPacket() *ArtCommandPacket {
	return &ArtCommandPacket{}
}

// MarshalBinary marshals an ArtCommandPacket into a byte slice.
func (p *ArtCommandPacket) MarshalBinary() ([]byte, error) {
	p.finish()
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtCommandPacket.
//TODO
func (p *ArtCommandPacket) UnmarshalBinary(b []byte) error {
	return p.validate()
}

// validate is used to validate the Packet.
func (p *ArtCommandPacket) validate() error {
	if p.OpCode != code.OpCommand {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtCommandPacket) finish() {
	p.OpCode = code.OpCommand
	p.version = version.Bytes()
	p.id = ArtNet
	p.estamanufacturer = [2]byte{0xff, 0xff}
}
