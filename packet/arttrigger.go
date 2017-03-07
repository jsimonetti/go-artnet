package packet

import (
	"bytes"
	"encoding/binary"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

var _ ArtNetPacket = &ArtTriggerPacket{}

// ArtTriggerPacket contains an ArtTrigger Packet.
//
// The ArtTrigger packet is used to send trigger macros to the network. The most common
// implementation involves a single controller broadcasting to all other devices.
// In some circumstances a controller may only wish to trigger a single device or a small
// group in which case unicast would be used
//
// The Key is an 8-bit number which defines the purpose of the packet. The interpretation
// of this field is dependent upon the Oem field. If the Oem field is set to a value other than
// 0xffff then the Key and SubKey fields are manufacturer specific
//
// The SubKey is an 8-bit number. The interpretation of this field is dependent upon the Oem field.
// If the Oem field is set to a value other than 0xffff then the Key and SubKey fields are
// manufacturer specific.
//
// The Payload is a fixed length array of 512, 8-bit bytes. The interpretation of this field is
// dependent upon the Oem field. If the Oem field is set to a value other than 0xffff then
// the Payload is manufacturer specific.
//
// However, when the Oem field = 0xffff the meaning of the Key, SubKey and Payload is:
//
// Key - Name     - SubKey:
//   0 - KeyAscii - The SubKey field contains an ASCII character which the receiving device should
//                  process as if it were a keyboard press. (Payload not used).
//   1 - KeyMacro - The SubKey field contains the number of a Macro which the receiving device
//                  should execute. (Payload not used).
//   2 - KeySoft  - The SubKey field contains a soft-key number which the receiving device should
//                  process as if it were a soft-key keyboard press. (Payload not used).
//   3 - KeyShow  - The SubKey field contains the number of a Show which the receiving device
//                  should run. (Payload not used).
//   4 - - 255      Undefined
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
type ArtTriggerPacket struct {
	// Inherit the Header header
	Header

	// this packet type contains a version
	version [2]byte

	// filler bytes
	filler [2]byte

	// Oem word describes the equipment manufacturer code of nodes that shall accept this trigger
	Oem uint16

	// Key is the Trigger Key
	Key uint8

	// SubKey is the Trigger SubKey
	SubKey uint8

	// Data interpretation of the payload is defined by the Key
	Data [512]byte
}

// NewArtTriggerPacket returns an ArtNetPacket with the correct OpCode
func NewArtTriggerPacket() *ArtTriggerPacket {
	return &ArtTriggerPacket{
		Header: Header{
			OpCode: code.OpTrigger,
			id:     ArtNet,
		},
		version: version.Bytes(),
	}
}

// MarshalBinary marshals an ArtTriggerPacket into a byte slice.
func (p *ArtTriggerPacket) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), p.validate()
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtTriggerPacket.
//TODO
func (p *ArtTriggerPacket) UnmarshalBinary(b []byte) error {
	return p.validate()
}

// artPacket is an empty method to sattisfy the ArtNetPacket interface.
func (p *ArtTriggerPacket) validate() error {
	if p.OpCode != code.OpTrigger {
		return errInvalidOpCode
	}
	return nil
}
