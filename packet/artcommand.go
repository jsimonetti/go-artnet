package packet

import (
	"github.com/jsimonetti/go-artnet/packet/code"
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
//	Command=Data&
//
// The ampersand is a break between commands. Also note that the text is capitalised for
// readability; it is case insensitive. Thus far, two commands are defined by Art-Net. It is
// anticipated that additional commands will be added as other manufacturers register commands
// which have industry wide relevance. These commands shall be transmitted with EstaMan = 0xFFFF.
//
// SwoutText - This command is used to re-programme the label associated with the
//
//	ArtPollReply->Swout fields. Syntax: "SwoutText=Playback&"
//
// SwinText  - This command is used to re-programme the label associated with the
//
//	ArtPollReply->Swin fields. Syntax: "SwinText=Record&"
//
// Packet Strategy:
//
//	Controller -  Receive:            Application Specific
//	              Unicast Transmit:   Application Specific
//	              Broadcast Transmit: Application Specific
//	Node -        Receive:            Application Specific
//	              Unicast Transmit:   Application Specific
//	              Broadcast Transmit: Application Specific
//	MediaServer - Receive:            Application Specific
//	              Unicast Transmit:   Application Specific
//	              Broadcast Transmit: Application Specific
type ArtCommandPacket struct {
	// Inherit the Header header
	Header

	// ESTAmanufacturer contains a code used to represent equipment manufacturer.
	ESTAmanufacturer [2]byte

	// Length indicates the length of the data
	Length uint16

	// Data is an ASCII string, null terminated. Max length is 512 bytes including the null terminator
	Data [512]byte
}

// NewArtCommandPacket returns an ArtNetPacket with the correct OpCode
func NewArtCommandPacket() *ArtCommandPacket {
	return &ArtCommandPacket{
		Header:           NewHeader(code.OpCommand),
		ESTAmanufacturer: [2]byte{0xff, 0xff},
	}
}

// MarshalBinary marshals an ArtCommandPacket into a byte slice.
func (p *ArtCommandPacket) MarshalBinary() ([]byte, error) {
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtCommandPacket.
func (p *ArtCommandPacket) UnmarshalBinary(b []byte) error {
	err := unmarshalPacket(p, b)
	if err != nil {
		return err
	}

	return p.Header.validate(code.OpCommand)
}
