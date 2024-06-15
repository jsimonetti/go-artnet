package artnet

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/artnettypes"
	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

// InputPort contains information for an input port
type InputPort struct {
	Address artnettypes.Address
	Type    code.PortType
	Status  code.GoodInput
}

// OutputPort contains information for an input port
type OutputPort struct {
	Address artnettypes.Address
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
	BindIndex artnettypes.BindIndex
	Port      uint16

	Report  code.NodeReport
	Status1 code.Status1
	Status2 code.Status2
	Status3 code.Status3

	BaseAddress artnettypes.Address
	InputPorts  []InputPort
	OutputPorts []OutputPort
}

// buildArtPollReply will return a ArtPollReplyPacket from the NodeConfig
// TODO: make this a more complete packet by adding the other NodeConfig fields
func (c NodeConfig) buildArtPollReply() *packet.ArtPollReplyPacket {
	p := packet.NewArtPollReplyPacket()

	p.Oem = c.OEM
	p.VersionInfo = c.Version
	p.UBEAVersion = c.BiosVersion
	p.Style = c.Type
	p.BindIndex = c.BindIndex
	p.Status1 = c.Status1
	p.Status2 = c.Status2
	p.Status3 = c.Status3
	p.NetSwitch = c.BaseAddress.Net
	p.SubSwitch = c.BaseAddress.SubUni
	p.NumPorts = c.NumberOfPorts()
	p.PortTypes = c.PortTypes()

	copy(p.IPAddress[0:4], c.IP.To4())
	copy(p.ESTAmanufacturer[0:2], c.Manufacturer)
	copy(p.ShortName[0:18], c.Name)
	copy(p.LongName[0:64], c.Description)
	copy(p.NodeReport[0:64], c.Report[0:64])
	copy(p.BindIP[0:4], c.BindIP.To4())
	copy(p.Macaddress[0:6], c.Ethernet)

	return p
}

// NumberOfPorts returns the count of node ports. This method assumes that
// NodeConfig is validated.
func (c NodeConfig) NumberOfPorts() uint16 {
	if len(c.InputPorts) > len(c.OutputPorts) {
		return uint16(len(c.InputPorts))
	}
	return uint16(len(c.OutputPorts))
}

// PortType merges the InputPorts and OutputPorts config into a single PortTypes
// definition. This method assumes that NodeConfig is validated.
func (c NodeConfig) PortTypes() [4]code.PortType {
	rsl := [4]code.PortType{}
	for i := 0; i < 4; i++ {
		tmp := code.PortType(0)
		if len(c.InputPorts) > i {
			tmp = tmp.WithInput(true)
			rsl[i] = tmp.WithType(c.InputPorts[i].Type.Type())
		}
		if len(c.OutputPorts) > i {
			tmp = tmp.WithOutput(true)
			rsl[i] = tmp.WithType(c.OutputPorts[i].Type.Type())
		}
	}
	return rsl
}

// validate will check the config and return an error if something is not valid.
// The main objective of this method is to check if the in- and output-ports configured
// by the user can be announced on the Art-Net network. It checks:
//
//   - At max 4 in- and/or outputs are supported per node.
//   - If a port supports in- and output at the same time the protocol type has to be
//     the same for the input port of the same index as the output port.
func (c NodeConfig) validate() error {
	if len(c.InputPorts) > 4 {
		return fmt.Errorf("validation error: more than 4 input ports configured (%d) for the node, this isn't supported by the library", len(c.InputPorts))
	}
	if len(c.OutputPorts) > 4 {
		return fmt.Errorf("validation error: more than 4 output ports configured (%d) for the node, this isn't supported by the library", len(c.InputPorts))
	}
	for i := 0; i < 4; i++ {
		if len(c.InputPorts) <= i || len(c.OutputPorts) <= i {
			continue
		}
		if c.InputPorts[i].Type.Type() != c.OutputPorts[i].Type.Type() {
			return fmt.Errorf(
				"validation error: the type (%s) of input port %d has a different type (%s) than output port %d, input and output ports with the same index must have the same type",
				c.InputPorts[i].Type.Type(),
				i+1,
				c.OutputPorts[i].Type.Type(),
				i+1,
			)
		}
	}
	return nil
}

// newNodeConfigFrom will return a Config from the information in the ArtPollReplyPacket
func newNodeConfigFrom(p *packet.ArtPollReplyPacket) NodeConfig {
	nodeConfig := NodeConfig{
		OEM:          p.Oem,
		Version:      p.VersionInfo,
		BiosVersion:  p.UBEAVersion,
		Manufacturer: decodeString(p.ESTAmanufacturer[:]),
		Type:         p.Style,
		Name:         decodeString(p.ShortName[:]),
		Description:  decodeString(p.LongName[:]),
		Report:       p.NodeReport,
		Ethernet:     p.Macaddress[:],
		IP:           p.IPAddress[:],
		BindIP:       p.BindIP[:],
		BindIndex:    p.BindIndex,
		Port:         p.Port,
		Status1:      p.Status1,
		Status2:      p.Status2,
		Status3:      p.Status3,
		BaseAddress: artnettypes.Address{
			Net:    p.NetSwitch,
			SubUni: p.SubSwitch << 4,
		},
	}

	for i := 0; i < int(p.NumPorts) && i < 4; i++ {
		if p.PortTypes[i].Output() {
			nodeConfig.OutputPorts = append(nodeConfig.OutputPorts, OutputPort{
				Address: artnettypes.Address{
					Net:    nodeConfig.BaseAddress.Net,
					SubUni: nodeConfig.BaseAddress.SubUni | p.SwOut[i],
				},
				Type:   p.PortTypes[i],
				Status: p.GoodOutput[i],
			})
		}
		if p.PortTypes[i].Input() {
			nodeConfig.InputPorts = append(nodeConfig.InputPorts, InputPort{
				Address: artnettypes.Address{
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

// decodeString will take a byte slice and create an ASCII string
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
