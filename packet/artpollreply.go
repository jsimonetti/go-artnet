package packet

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/packet/code"
)

var _ ArtNetPacket = &ArtPollReplyPacket{}

// ArtPollReplyPacket contains an ArtPollReply Packet.
type ArtPollReplyPacket struct {
	// Inherit the Header header
	Header

	// IPAddress is the Node’s IPv4 address. When binding is implemented, bound nodes may
	// share the root node’s IP Address and the BindIndex is used to differentiate the nodes.
	IPAddress net.IP

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
	ESTAmanufacturer string

	// ShortName for the Node. The Controller uses the ArtAddress packet to program this
	// string. Max length is 17 characters. This is a fixed length field, although the string
	// it contains can be shorter than the field.
	ShortName string

	// LongName for the Node. The Controller uses the ArtAddress packet to program this string.
	// Max length is 63. This is a fixed length field, although the string it contains can be
	// shorter than the field.
	LongName string

	// NodeReport is a textual report of the Node’s operating status or operational errors.
	// It is primarily intended for ‘engineering’ data.
	NodeReport [64]code.NodeReportCode

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

	// TODO (4x uint8, array)
	SwIn [4]uint8

	// TODO (4x uint8, array)
	SwOut [4]uint8

	// SwVideo is set to 00 when video display is showing local data. Set to 01 when video
	// is showing ethernet data. The field is now deprecated
	SwVideo uint8

	// SwMacro shows if the Node supports macro key inputs, this byte represents the trigger values.
	SwMacro uint8

	// SwRemote show if the Node supports remote trigger inputs, this byte represents the trigger values.
	SwRemote uint8

	// spare bytes
	spare [3]byte

	// Style code defines the equipment style of the device.
	Style code.StyleCode

	// Macaddress of the Node. Set to zero if node cannot supply this information.
	Macaddress net.HardwareAddr

	// BindIP is the IP of the root device if this unit is part of a larger or modular product.
	BindIP net.IP

	// BindIndex represents the order of bound devices. A lower number means closer to root device.
	// A value of 1 means root device.
	BindIndex uint8

	// Status2 indicates Product capabilities
	Status2 code.Status2

	// filler bytes. Transmit as zero. For future expansion.
	filler [26]byte
}

// NewArtPollReplyPacket returns a new ArtPollReply Packet
func NewArtPollReplyPacket() *ArtPollReplyPacket {
	return &ArtPollReplyPacket{
		Header: Header{
			OpCode: code.OpPollReply,
			id:     ArtNet,
		},
	}
}

// MarshalBinary marshals an ArtPollReplyPacket into a byte slice.
// TODO
func (p *ArtPollReplyPacket) MarshalBinary() ([]byte, error) {
	return nil, p.validate()
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollReplyPacket.
func (p *ArtPollReplyPacket) UnmarshalBinary(b []byte) error {
	if err := p.Header.unmarshal(b[:10]); err != nil {
		return err
	}
	if len(b) != 238 {
		return fmt.Errorf("invalid packet length received. want: 238, got: %d", len(b))
	}

	p.IPAddress = b[10:14]
	p.Port = uint16(b[14]) | uint16(b[15])<<8
	p.VersionInfo = uint16(b[16]) | uint16(b[17])<<8

	p.NetSwitch = b[18]
	p.SubSwitch = b[19]
	p.Oem = uint16(b[20]) | uint16(b[21])<<8

	p.UBEAVersion = b[22]
	p.Status1 = code.Status1(b[23])

	man := []byte{uint8(b[24] & 0xff), uint8(uint16(b[25]) << 8)}
	for _, c := range man {
		if c == 0x00 {
			break
		}
		p.ESTAmanufacturer += string(c)
	}

	for _, c := range b[26:44] {
		if c == 0x00 {
			break
		}
		p.ShortName += string(c)
	}

	for _, c := range b[44:108] {
		if c == 0x00 {
			break
		}
		p.LongName += string(c)
	}

	for i, r := range b[108:172] {
		p.NodeReport[i] = code.NodeReportCode(r)
	}

	p.NumPorts = uint16(b[172]) | uint16(b[173])<<8
	for i, r := range b[174:178] {
		p.PortTypes[i] = code.PortType(r)
	}

	for i, r := range b[178:182] {
		p.GoodInput[i] = code.GoodInput(r)
	}

	for i, r := range b[182:186] {
		p.GoodOutput[i] = code.GoodOutput(r)
	}

	for i, r := range b[186:190] {
		p.SwIn[i] = r
	}

	for i, r := range b[190:194] {
		p.SwOut[i] = r
	}
	p.SwVideo = b[194]
	p.SwMacro = b[195]
	p.SwRemote = b[196]
	p.spare = [3]byte{b[197], b[198], b[199]}
	p.Style = code.StyleCode(b[200])

	p.Macaddress = b[200:206]
	p.BindIP = b[206:210]
	p.BindIndex = b[210]
	p.Status2 = code.Status2(b[211])

	for i := 0; i < 26; i++ {
		p.filler[i] = b[212+i]
	}

	return p.validate()
}

// artPacket is an empty method to sattisfy the ArtNetPacket interface.
func (p *ArtPollReplyPacket) validate() error {
	if p.OpCode != code.OpPollReply {
		return errInvalidOpCode
	}
	if !code.ValidStyle(p.Style) {
		return errInvalidStyle
	}
	return nil
}
