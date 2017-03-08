package packet

import (
	"github.com/jsimonetti/go-artnet/packet/code"
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

	// Filler bytes
	_ [2]byte

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
	return &ArtTriggerPacket{}
}

// MarshalBinary marshals an ArtTriggerPacket into a byte slice.
func (p *ArtTriggerPacket) MarshalBinary() ([]byte, error) {
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtTriggerPacket.
func (p *ArtTriggerPacket) UnmarshalBinary(b []byte) error {
	return unmarshalPacket(p, b)
}

// validate is used to validate the Packet.
func (p *ArtTriggerPacket) validate() error {
	if err := p.Header.validate(); err != nil {
		return err
	}
	if p.OpCode != code.OpTrigger {
		return errInvalidOpCode
	}
	return nil
}

// finish is used to finish the Packet for sending.
func (p *ArtTriggerPacket) finish() {
	p.Header.finish()
}
