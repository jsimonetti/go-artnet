package artnet

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
