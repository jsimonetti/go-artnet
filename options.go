package artnet

import (
	"net"

	"github.com/jsimonetti/go-artnet/packet"
)

// Option is a functional option handler for Controller.
type Option func(*Controller) error

// SetOption runs a functional option against Controller.
func (c *Controller) SetOption(option Option) error {
	return option(c)
}

// MaxFPS sets the maximum amount of updates sent out per second
func MaxFPS(fps int) Option {
	return func(c *Controller) error {
		c.maxFPS = fps
		return nil
	}
}

// BroadcastAddr sets the broadcast address to use; defaults to 2.255.255.255:6454
func BroadcastAddr(addr net.UDPAddr) Option {
	return func(c *Controller) error {
		c.broadcastAddr = addr
		return nil
	}
}

// ListenAddr sets the listen address and port to use; defaults to :6454 if unset
func ListenAddress(addr net.UDPAddr) Option {
	return func(c *Controller) error {
		NodeListenAddress(addr)(c.cNode)
		return nil
	}
}

// ListenIP sets the listen IP to use; defaults to :6454 if unset
func ListenIP(ip net.IP) Option {
	return func(c *Controller) error {
		NodeListenIP(ip)(c.cNode)
		return nil
	}
}

// NodeOption is a functional option handler for Node.
type NodeOption func(*Node) error

// SetOption runs a functional option against Node.
func (n *Node) SetOption(option NodeOption) error {
	return option(n)
}

// NodeBroadcastAddress sets the broadcast address to use; defaults to 2.255.255.255:6454
func NodeBroadcastAddress(addr net.UDPAddr) NodeOption {
	return func(n *Node) error {
		n.broadcastAddr = addr
		return nil
	}
}

// NodeListenAddress sets the listen address and port to use; defaults to :6454 if unset
func NodeListenAddress(addr net.UDPAddr) NodeOption {
	return func(n *Node) error {
		n.listenAddr = addr
		return nil
	}
}

// NodeListenIP sets the listen IP to use; defaults to :6454 if unset
func NodeListenIP(ip net.IP) NodeOption {
	return func(n *Node) error {
		n.listenAddr = net.UDPAddr{IP: ip, Port: packet.ArtNetPort}
		return nil
	}
}
