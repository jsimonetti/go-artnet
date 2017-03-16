package packet

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/packet/code"
)

var _ ArtNetPacket = &ArtPollReplyPacket{}

// ArtPollReplyPacket contains an ArtPollReply Packet.
//
// A device, in response to a Controller’s ArtPoll, sends the ArtPollReply. This packet
// is also broadcast to the Directed Broadcast address by all Art-Net devices on power up.
//
// Packet Strategy:
//  All devices - Receive:            No Art-Net action.
//                Unicast Transmit:   Not Allowed.
//                Broadcast Transmit: Directed Broadcasts this packet in response to an ArtPoll.
type ArtPollReplyPacket struct {
	// ID is an Array of 8 characters, the final character is a null termination.
	// Value should be []byte{‘A’,‘r’,‘t’,‘-‘,‘N’,‘e’,‘t’,0x00}
	// ArtPollReply is the only packet not containing version, so do this here
	ID [8]byte

	// OpCode defines the class of data following within this UDP packet.
	// Transmitted low byte first.
	OpCode code.OpCode

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

	// BindIP is the IP of the root device if this unit is part of a larger or modular product.
	BindIP [4]byte

	// BindIndex represents the order of bound devices. A lower number means closer to root device.
	// A value of 1 means root device.
	BindIndex uint8

	// Status2 indicates Product capabilities
	Status2 code.Status2

	// Filler bytes. Transmit as zero. For future expansion.
	_ [26]byte
}

// NewArtPollReplyPacket returns a new ArtPollReply Packet
func NewArtPollReplyPacket() *ArtPollReplyPacket {
	return &ArtPollReplyPacket{}
}

// MarshalBinary marshals an ArtPollReplyPacket into a byte slice.
func (p *ArtPollReplyPacket) MarshalBinary() ([]byte, error) {
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtPollReplyPacket.
func (p *ArtPollReplyPacket) UnmarshalBinary(b []byte) error {
	return unmarshalPacket(p, b)
}

// validate is used to validate the Packet.
func (p *ArtPollReplyPacket) validate() error {
	// swap endianness
	p.OpCode = code.OpCode(swapUint16(uint16(p.OpCode)))
	if p.OpCode != code.OpPollReply {
		return errInvalidOpCode
	}
	p.Port = swapUint16(p.Port)

	// It appears not all software sends the port low byte first
	if p.Port != ArtNetPort {
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

// finish is used to finish the Packet for sending.
func (p *ArtPollReplyPacket) finish() {
	p.ID = ArtNet
	p.OpCode = code.OpCode(swapUint16(uint16(code.OpPollReply)))
	p.Port = swapUint16(p.Port)
}

// Address contains a universe address
type Address struct {
	Net    uint8 // 0-128
	SubUni uint8
}

// String returns a string representation of Address
func (a Address) String() string {
	return fmt.Sprintf("%d:%d.%d", a.Net, (a.SubUni >> 4), a.SubUni&0x0f)
}

// Integer returns the integer representation of Address
func (a Address) Integer() int {
	return int(uint16(a.Net)<<8 | uint16(a.SubUni))
}

// InputPort contains information for an input port
type InputPort struct {
	Address Address
	Type    code.PortType
	Status  code.GoodInput
}

// OutputPort contains information for an input port
type OutputPort struct {
	Address Address
	Type    code.PortType
	Status  code.GoodOutput
}

// NodeConfig is a representation of a single node.
type NodeConfig struct {
	OEM          uint16
	Version      uint16
	BiosVersion  uint8
	Manufacturer string
	Type         string
	Name         string
	Description  string
	Report       []code.NodeReportCode
	Ethernet     net.HardwareAddr
	IP           net.IP
	BindIP       net.IP
	BindIndex    uint8
	Port         uint16
	Status1      code.Status1
	Status2      code.Status2

	BaseAddress Address
	InputPorts  []InputPort
	OutputPorts []OutputPort
}

// NodeConfig returns a NodeConfig based on the information in the packet
func (p *ArtPollReplyPacket) NodeConfig() NodeConfig {

	nodeConfig := NodeConfig{
		OEM:          p.Oem,
		Version:      p.VersionInfo,
		BiosVersion:  p.UBEAVersion,
		Manufacturer: decodeString(p.ESTAmanufacturer[:]),
		Type:         p.Style.String(),
		Name:         decodeString(p.ShortName[:]),
		Description:  decodeString(p.LongName[:]),
		Report:       p.NodeReport[:],
		Ethernet:     p.Macaddress[:],
		IP:           p.IPAddress[:],
		BindIP:       p.BindIP[:],
		BindIndex:    p.BindIndex,
		Port:         p.Port,
		Status1:      p.Status1,
		Status2:      p.Status2,
		BaseAddress: Address{
			Net:    p.NetSwitch,
			SubUni: p.SubSwitch,
		},
	}

	for i := 0; i < int(p.NumPorts) && i < 4; i++ {
		if p.PortTypes[i].Output() {
			nodeConfig.OutputPorts = append(nodeConfig.OutputPorts, OutputPort{
				Address: Address{
					Net:    nodeConfig.BaseAddress.Net,
					SubUni: nodeConfig.BaseAddress.SubUni | p.SwOut[i],
				},
				Type:   p.PortTypes[i],
				Status: p.GoodOutput[i],
			})
		}
		if p.PortTypes[i].Input() {
			nodeConfig.InputPorts = append(nodeConfig.InputPorts, InputPort{
				Address: Address{
					Net:    nodeConfig.BaseAddress.Net,
					SubUni: nodeConfig.BaseAddress.SubUni | p.SwIn[i],
				},
				Type:   p.PortTypes[i],
				Status: p.GoodInput[i],
			})
		}
	}

	return nodeConfig
}

func decodeString(b []byte) string {
	var str string
	for _, c := range b {
		if c != 0 {
			str += string(c)
		}
	}
	return str
}
