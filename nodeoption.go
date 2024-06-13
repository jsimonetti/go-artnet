package artnet

import (
	"net"

	"github.com/jsimonetti/go-artnet/packet/code"
)

// NodeOption is a functional option handler for Node.
type NodeOption func(*Node) error

// SetOption runs a functional option against Node.
func (n *Node) SetOption(option NodeOption) error {
	return option(n)
}

// NodeBroadcastAddress sets the broadcast address to use; defaults to 2.255.255.255:6454
func NodeOptionBroadcastAddress(addr net.UDPAddr) NodeOption {
	return func(n *Node) error {
		n.broadcastAddr = addr
		return nil
	}
}

// NodeOptionPacketHandler sets the packetHandler for a given OpCode
func NodeOptionPacketHandler(opcode code.OpCode, handler PacketHandler) NodeOption {
	return func(n *Node) error {
		n.packetHandlers[opcode] = handler
		return nil
	}

}
