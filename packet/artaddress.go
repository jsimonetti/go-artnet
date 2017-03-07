package packet

import (
	"bytes"
	"encoding/binary"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

var _ ArtNetPacket = &ArtAddressPacket{}

// ArtAddressPacket contains an ArtAddress Packet.
//
// A Controller or monitoring device on the network can reprogram numerous controls of a
// node remotely. This, for example, would allow the lighting console to re-route DMX512
// data at remote locations. This is achieved by sending an ArtAddress packet to the Node’s
// IP address. (The IP address is returned in the ArtPoll packet). The node replies with an
// ArtPollReply packet.
// Fields 5 to 13 contain the data that will be programmed into the node
//
// Packet Strategy:
//  Controller -  Receive:            No Action
//                Unicast Transmit:   Controller transmits to a specific node IP address
//                Broadcast Transmit: Not Allowed
//  Node -        Receive:            Reply by broadcasting ArtPollReply
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Not Allowed
//  MediaServer - Receive:            Reply by broadcasting ArtPollReply
//                Unicast Transmit:   Not Allowed
//                Broadcast Transmit: Not Allowed
type ArtAddressPacket struct {
	// Inherit the Header header
	Header

	// this packet type contains a version
	version [2]byte

	// NetSwitch contains Bits 14-8 of the 15 bit Port-Address are encoded into the bottom 7
	// bits of this field. This is used in combination with SubSwitch and SwIn[] or SwOut[]
	// to produce the full universe address.
	NetSwitch uint8

	// BindIndex defines the bound node which originated this packet and is used to uniquely identify
	// the bound node when identical IP addresses are in use. This number represents the order of bound
	// devices. A lower number means closer to root device. A value of 1 means root device.
	BindIndex uint8

	// ShortName for the Node. The Controller uses the ArtAddress packet to program this
	// string. Max length is 17 characters. This is a fixed length field, although the string
	// it contains can be shorter than the field.
	ShortName string

	// LongName for the Node. The Controller uses the ArtAddress packet to program this string.
	// Max length is 63. This is a fixed length field, although the string it contains can be
	// shorter than the field.
	LongName string

	// SwIn Bits 3-0 of the 15 bit Port-Address for each of the 4
	// possible input ports are encoded into the low nibble
	SwIn [4]uint8

	// SwOut Bits 3-0 of the 15 bit Port-Address for each of the 4
	// possible output ports are encoded into the low nibble.
	SwOut [4]uint8

	// SubSwitch contains Bits 7-4 of the 15 bit Port-Address are encoded into the bottom 4
	// bits of this field. This is used in combination with NetSwitch and SwIn[] or SwOut[]
	// to produce the full universe address.
	SubSwitch uint8

	// SwVideo is set to 00 when video display is showing local data. Set to 01 when video
	// is showing ethernet data. The field is now deprecated
	SwVideo uint8

	// Command contains Node configuration commands. Note that Ltp / Htp settings should be
	// retained by the node during power cycling
	Command uint8
}

// NewArtAddressPacket returns an ArtNetPacket with the correct OpCode
func NewArtAddressPacket() *ArtAddressPacket {
	return &ArtAddressPacket{}
}

// MarshalBinary marshals an ArtAddressPacket into a byte slice.
func (p *ArtAddressPacket) MarshalBinary() ([]byte, error) {
	p.finish()
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtAddressPacket.
//TODO
func (p *ArtAddressPacket) UnmarshalBinary(b []byte) error {
	return p.validate()
}

// validate is used to validate the Packet.
func (p *ArtAddressPacket) validate() error {
	if p.OpCode != code.OpAddress {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtAddressPacket) finish() {
	p.OpCode = code.OpAddress
	p.version = version.Bytes()
	p.id = ArtNet
}
