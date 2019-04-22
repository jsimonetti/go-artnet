package artnet

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

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
	Type         code.StyleCode
	Name         string
	Description  string

	Ethernet  net.HardwareAddr
	IP        net.IP
	BindIP    net.IP
	BindIndex uint8
	Port      uint16

	Report  []code.NodeReportCode
	Status1 code.Status1
	Status2 code.Status2

	BaseAddress Address
	InputPorts  []InputPort
	OutputPorts []OutputPort
}

// ArtPollReplyFromConfig will return a ArtPollReplyPacket from the NodeConfig
// TODO: make this a more complete packet by adding the other NodeConfig fields
func ArtPollReplyFromConfig(c NodeConfig) *packet.ArtPollReplyPacket {
	// create an ArtPollReply packet to send out with the ArtPoll packet
	p := &packet.ArtPollReplyPacket{
		OpCode:      code.OpPollReply,
		Port:        c.Port,
		Oem:         c.OEM,
		VersionInfo: c.Version,
		UBEAVersion: c.BiosVersion,
		Style:       c.Type,
		BindIndex:   c.BindIndex,
		Status1:     c.Status1,
		Status2:     c.Status2,
		NetSwitch:   c.BaseAddress.Net,
		SubSwitch:   c.BaseAddress.SubUni,
	}
	copy(p.IPAddress[0:4], c.IP.To4())
	copy(p.ESTAmanufacturer[0:2], c.Manufacturer)
	copy(p.ShortName[0:18], c.Name)
	copy(p.LongName[0:64], c.Description)
	copy(p.NodeReport[0:64], c.Report)
	copy(p.BindIP[0:4], c.BindIP.To4())
	copy(p.Macaddress[0:6], c.Ethernet)

	return p
}

// ConfigFromArtPollReply will return a Config from the information in the ArtPollReplyPacket
func ConfigFromArtPollReply(p packet.ArtPollReplyPacket) NodeConfig {
	nodeConfig := NodeConfig{
		OEM:          p.Oem,
		Version:      p.VersionInfo,
		BiosVersion:  p.UBEAVersion,
		Manufacturer: decodeString(p.ESTAmanufacturer[:]),
		Type:         p.Style,
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

// decodeString will take a byteslice and create an ASCII string
// the ASCII strings are 0 terminated
func decodeString(b []byte) (str string) {
	for _, c := range b {
		if c == 0 {
			return
		}
		str += string(c)
	}
	return
}
