package code

// Status2 indicates Product capabilities
type Status2 uint8

// WithBrowser sets if product supports web browser configuration
func (s Status2) WithBrowser(enable bool) Status2 {

	if enable {
		return s | (1 << 0)
	}
	return s | (0 << 0)
}

// Browser indicates if product supports web browser configuration
func (s Status2) Browser() bool {
	return s&(1<<0) > 0
}

// WithDHCP sets if product IP is DHCP configured
func (s Status2) WithDHCP(enable bool) Status2 {

	if enable {
		return s | (1 << 1)
	}
	return s | (0 << 1)
}

// DHCP indicates if product IP is DHCP configured
func (s Status2) DHCP() bool {
	return s&(1<<1) > 0
}

// WithDHCPCapable sets if product is capable of DHCP
func (s Status2) WithDHCPCapable(enable bool) Status2 {

	if enable {
		return s | (1 << 2)
	}
	return s | (0 << 2)
}

// DHCPCapable indicates if product is capable of DHCP
func (s Status2) DHCPCapable() bool {
	return s&(1<<2) > 0
}

// WithPort15 sets if product supports 15 bit Port-Address (Art-Net 3 or 4)
func (s Status2) WithPort15(enable bool) Status2 {

	if enable {
		return s | (1 << 3)
	}
	return s | (0 << 3)
}

// Port15 indicates if product supports 15 bit Port-Address (Art-Net 3 or 4)
func (s Status2) Port15() bool {
	return s&(1<<3) > 0
}

// WithSwitch sets if product is able to switch between Art-Net and sACN
func (s Status2) WithSwitch(enable bool) Status2 {

	if enable {
		return s | (1 << 4)
	}
	return s | (0 << 4)
}

// Switch indicates if product is able to switch between Art-Net and sACN
func (s Status2) Switch() bool {
	return s&(1<<4) > 0
}

// WithSquawk sets if product is squawking
func (s Status2) WithSquawk(enable bool) Status2 {

	if enable {
		return s | (1 << 5)
	}
	return s | (0 << 5)
}

// Squawk indicates if product is squawking
func (s Status2) Squawk() bool {
	return s&(1<<5) > 0
}

// String returns a string representation of TalkToMe
func (s Status2) String() string {
	browser, dhcp, dhcpcap, port15, swtch, squawk := "no", "no", "no", "no", "no", "no"
	if s.Browser() {
		browser = "yes"
	}
	if s.DHCP() {
		dhcp = "yes"
	}
	if s.DHCPCapable() {
		dhcpcap = "yes"
	}
	if s.Port15() {
		port15 = "yes"
	}
	if s.Switch() {
		swtch = "yes"
	}
	if s.Squawk() {
		squawk = "yes"
	}

	return "Status2: Browser: " + browser + ", DHCP: " + dhcp + ", DHCPCapable: " + dhcpcap + ", Port 15bit: " + port15 + ", CanSwitch: " + swtch + ", Squawking: " + squawk
}
