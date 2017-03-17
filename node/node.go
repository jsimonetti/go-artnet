package node

import (
	"fmt"
	"net"

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
	Type         string
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
