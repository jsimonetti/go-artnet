package packet

import (
	"fmt"

	"github.com/jsimonetti/go-artnet/artnettypes"
	"github.com/jsimonetti/go-artnet/packet/code"
)

var _ ArtNetPacket = &ArtPollReplyPacket{}

// ArtPollReplyPacket contains an ArtPollReply Packet.
//
// A device, in response to a Controller’s ArtPoll, sends the ArtPollReply. This packet
// is also broadcast to the Directed Broadcast address by all Art-Net devices on power up.
//
// Packet Strategy:
//
//	All devices - Receive:            No Art-Net action.
//	              Unicast Transmit:   Not Allowed.
//	              Broadcast Transmit: Directed Broadcasts this packet in response to an ArtPoll.
type ArtPollReplyPacket struct {
	// Inherit the HeaderWithoutVersion
	HeaderWithoutVersion

	// IPAddress is the Node’s IPv4 address. When binding is implemented, bound nodes may
	// share the root node’s IP Address and the BindIndex is used to differentiate the nodes.
	IPAddress [4]byte

	// Port is always 0x1936 Transmitted low byte first.
	Port uint16

	// VersionInfo contains the Node’s firmware revision number. The Controller should only
	// use this field to decide if a firmware update should proceed. The convention is that
	// a higher number is a more recent release of firmware.
	VersionInfo uint16

	// NetSwitch contains Bits 14-8 of the 15 bit Port-Address are encoded into the bottom 7
	// bits of this field. This is used in combination with SubSwitch and SwIn[] or SwOut[]
	// to produce the full universe address.
	NetSwitch uint8

	// SubSwitch contains Bits 7-4 of the 15 bit Port-Address are encoded into the bottom 4
	// bits of this field. This is used in combination with NetSwitch and SwIn[] or SwOut[]
	// to produce the full universe address.
	SubSwitch uint8

	// Oem word describes the equipment vendor and the feature set available.
	Oem uint16

	// UBEAVersion contains the firmware version of the User Bios Extension Area (UBEA).
	// If the UBEA is not programmed, this field contains zero.
	UBEAVersion uint8

	// Status1 indicates General Status register containing bit fields as follows.
	Status1 code.Status1

	// ESTAmanufacturer contains a code used to represent equipment manufacturer.
	// They are assigned by ESTA. This field can be interpreted as two ASCII bytes
	// representing the manufacturer initials.
	ESTAmanufacturer [2]byte

	// ShortName for the Node. The Controller uses the ArtAddress packet to program this
	// string. Max length is 17 characters. This is a fixed length field, although the string
	// it contains can be shorter than the field.
	ShortName [18]byte

	// LongName for the Node. The Controller uses the ArtAddress packet to program this string.
	// Max length is 63. This is a fixed length field, although the string it contains can be
	// shorter than the field.
	LongName [64]byte

	// NodeReport is a textual report of the Node’s operating status or operational errors.
	// It is primarily intended for ‘engineering’ data.
	NodeReport code.NodeReport

	// NumPorts describes the number of input or output ports. If number of inputs is not
	// equal to number of outputs, the largest value is taken. Zero is a legal value if no
	// input or output ports are implemented. The maximum value is 4. Nodes can ignore this
	// field as the information is implicit in PortTypes.
	NumPorts uint16

	// PortTypes defines the operation and protocol of each channel
	PortTypes [4]code.PortType

	// GoodInput defines input status of the node
	GoodInput [4]code.GoodInput

	// GoodOutput defines output status of the node
	GoodOutput [4]code.GoodOutput

	// SwIn Bits 3-0 of the 15 bit Port-Address for each of the 4
	// possible input ports are encoded into the low nibble
	SwIn [4]uint8

	// SwOut Bits 3-0 of the 15 bit Port-Address for each of the 4
	// possible output ports are encoded into the low nibble.
	SwOut [4]uint8

	// SwVideo is set to 00 when video display is showing local data. Set to 01 when video
	// is showing ethernet data. The field is now deprecated
	SwVideo uint8

	// SwMacro shows if the Node supports macro key inputs, this byte represents the trigger values.
	SwMacro code.SwMacro

	// SwRemote show if the Node supports remote trigger inputs, this byte represents the trigger values.
	SwRemote code.SwRemote

	// Spare bytes
	_ [3]byte

	// Style code defines the equipment style of the device.
	Style code.StyleCode

	// Macaddress of the Node. Set to zero if node cannot supply this information.
	Macaddress [6]byte

	// ============================================================================================
	//                                  Optional Fields - 32 Bytes
	//                   (only 207 bytes required for backwards compatibility)
	// ============================================================================================

	// BindIP is the IP of the root device if this unit is part of a larger or modular product.
	BindIP [4]byte

	// BindIndex represents the order of bound devices. A lower number means closer to root device.
	// A value of 1 means root device.
	BindIndex artnettypes.BindIndex

	// Status2 indicates Product capabilities
	Status2 code.Status2

	// GoodOutputB defines output status of the node
	GoodOutputB [4]byte

	// Status3 indicates general product state
	Status3 code.Status3

	// DefaultResponderUID
	DefaultResponderUID [6]byte

	// User is available for user Specific Data
	User [2]byte

	// RefreshRate specify the maximum refresh rate, expressed in Hz
	// Designed to allow refresh rates above DMX512 rate
	RefreshRate uint16

	// Filler bytes. Transmit as zero. For future expansion.
	_ [11]byte
}

const (
	minimumArtPollReplyPacketSize int = 207
	maximumArtPollReplyPacketSize int = minimumArtPollReplyPacketSize + 32
)

// NewArtPollReplyPacket returns a new ArtPollReply Packet
func NewArtPollReplyPacket() *ArtPollReplyPacket {
	return &ArtPollReplyPacket{
		HeaderWithoutVersion: NewHeaderWithoutVersion(code.OpPollReply),
		Port:                 ArtNetPort,
	}
}

// MarshalBinary marshals an ArtPollReplyPacket into a byte slice.
func (p *ArtPollReplyPacket) MarshalBinary() ([]byte, error) {
	p.Port = swapUint16(p.Port)
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollReplyPacket.
func (p *ArtPollReplyPacket) UnmarshalBinary(b []byte) error {
	err := checkPadAndUnmarshalPacket(p, b, minimumArtPollReplyPacketSize, maximumArtPollReplyPacketSize)
	if err != nil {
		return err
	}

	if err := p.HeaderWithoutVersion.validate(code.OpPollReply); err != nil {
		return err
	}

	p.Port = swapUint16(p.Port)

	// It appears not all software sends the port low byte first
	// so make an extra check here
	if p.Port != ArtNetPort {
		// since port could be swapped, swap representation
		p.Port = swapUint16(p.Port)
		if p.Port != ArtNetPort {
			return fmt.Errorf("invalid port: want: %d, got: %d", ArtNetPort, p.Port)
		}
	}

	if !code.ValidStyle(p.Style) {
		return errInvalidStyleCode
	}

	return nil
}
