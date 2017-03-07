package packet

import (
	"fmt"
	"net"

	"github.com/jsimonetti/artnet/packet/code"
)

var _ ArtNetPacket = &ArtPollReplyPacket{}

// ArtPollReplyPacket contains an ArtPollReply Packet.
type ArtPollReplyPacket struct {
	// Inherit the Packet header
	Packet

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

	// TODO
	Status1 uint8

	// ESTAmanufacturer contains a code used to represent equipment manufacturer.
	// They are assigned by ESTA. This field can be interpreted as two ASCII bytes
	// representing the manufacturer initials.
	ESTAmanufacturer uint16

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

	// TODO (4x uint8, array)
	PortTypes [4]uint8

	// TODO (4x uint8, array)
	GoodInput [4]uint8

	// TODO (4x uint8, array)
	GoodOutput [4]uint8

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
	Style uint8

	// Macaddress of the Node. Set to zero if node cannot supply this information.
	Macaddress net.HardwareAddr

	// BindIP is the IP of the root device if this unit is part of a larger or modular product.
	BindIP net.IP

	// BindIndex represents the order of bound devices. A lower number means closer to root device.
	// A value of 1 means root device.
	BindIndex uint8

	//TODO
	Status2 uint8

	// filler bytes. Transmit as zero. For future expansion.
	filler [26]byte
}

// NewArtPollReplyPacket returns a new ArtPollReply Packet
func NewArtPollReplyPacket() *ArtPollReplyPacket {
	return &ArtPollReplyPacket{}
}

// MarshalBinary marshals an ArtPollReplyPacket into a byte slice.
// TODO
func (p *ArtPollReplyPacket) MarshalBinary() ([]byte, error) {
	return nil, nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollReplyPacket.
func (p *ArtPollReplyPacket) UnmarshalBinary(b []byte) error {

	if len(b) != 228 {
		//header size is 10, so add that to error msg to prevent confusion
		return fmt.Errorf("invalid packet length received. want: 238, got: %d", len(b)+10)
	}

	p.IPAddress = b[0:4]
	p.Port = uint16(b[4]) | uint16(b[5])<<8
	p.VersionInfo = uint16(b[6]) | uint16(b[7])<<8

	p.NetSwitch = b[8]
	p.SubSwitch = b[9]
	p.Oem = uint16(b[10]) | uint16(b[11])<<8

	p.UBEAVersion = b[12]
	p.Status1 = b[13]
	p.ESTAmanufacturer = uint16(b[14]) | uint16(b[15])<<8

	for _, c := range b[16:34] {
		if c == 0x00 {
			break
		}
		p.ShortName += string(c)
	}

	for _, c := range b[34:98] {
		if c == 0x00 {
			break
		}
		p.LongName += string(c)
	}

	for i, r := range b[98:162] {
		p.NodeReport[i] = code.NodeReportCode(r)
	}

	p.NumPorts = uint16(b[162]) | uint16(b[163])<<8
	for i, r := range b[164:168] {
		p.PortTypes[i] = r
	}

	for i, r := range b[168:172] {
		p.GoodInput[i] = r
	}

	for i, r := range b[172:176] {
		p.GoodOutput[i] = r
	}

	for i, r := range b[176:180] {
		p.SwIn[i] = r
	}

	for i, r := range b[180:184] {
		p.SwOut[i] = r
	}
	p.SwVideo = b[184]
	p.SwMacro = b[185]
	p.SwRemote = b[186]
	p.spare = [3]byte{b[187], b[188], b[189]}
	p.Style = b[190]

	p.Macaddress = b[190:196]
	p.BindIP = b[196:200]
	p.BindIndex = b[200]
	p.Status2 = b[201]

	for i := 0; i < 26; i++ {
		p.filler[i] = b[202+i]
	}

	return nil
}

// TalkToMe sets the behaviour of a Node
// only bits 1-4 matter, rest is zero
type TalkToMe uint8

// WithReplyOnChange allows the Controller to be informed of changes
// without the need to continuously poll.
func (t TalkToMe) WithReplyOnChange(enable bool) TalkToMe {

	if enable {
		return t | (1 << 1)
	}
	return t | (0 << 1)
}

// ReplyOnChange returns the status of the bit 1
func (t TalkToMe) ReplyOnChange() bool {
	return t&(1<<1) > 0
}

// WithDiagnostics sends diagnostics messages
func (t TalkToMe) WithDiagnostics(enable bool) TalkToMe {
	if enable {
		return t | (1 << 2)
	}
	return t | (0 << 2)
}

// Diagnostics returns the stats of the bit 2
func (t TalkToMe) Diagnostics() bool {
	return t&(1<<2) > 0
}

// WithDiagUnicast determines wether diagnostics messages are unicast or broadcast
func (t TalkToMe) WithDiagUnicast(enable bool) TalkToMe {
	if enable {
		return t | (1 << 3)
	}
	return t | (0 << 3)
}

// DiagUnicast returns the stats of the bit 3
func (t TalkToMe) DiagUnicast() bool {
	return t&(1<<3) > 0
}

// WithVLC enable or disable VLC transmission
func (t TalkToMe) WithVLC(enable bool) TalkToMe {
	if enable {
		return t | (1 << 4)
	}
	return t | (0 << 4)
}

// VLC returns the stats of the bit 4
func (t TalkToMe) VLC() bool {
	return t&(1<<4) > 0
}

// String returns a string representation of TalkToMe
func (t TalkToMe) String() string {
	roc, diag, uni, vlc := "no", "no", "no", "no"
	if t.ReplyOnChange() {
		roc = "yes"
	}
	if t.Diagnostics() {
		diag = "yes"
	}
	if t.DiagUnicast() {
		uni = "yes"
	}
	if t.VLC() {
		vlc = "yes"
	}
	return "TalkToMe: ReplyOnChange: " + roc + ", Diagnostics: " + diag + " (as unicast: " + uni + "), VLC: " + vlc
}
