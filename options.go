package artnet

import (
	"net"
	"time"
)

// Option is a functional option handler for Controller.
type Option func(*Controller) error

// SetOption runs a functional option against Controller.
func (c *Controller) SetOption(option Option) error {
	return option(c)
}

// UpdateInterval sets the interval on which dmx updates will be sent out to known devices
func UpdateInterval(d time.Duration) Option {
	return func(c *Controller) error {
		c.updateInterval = d
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
