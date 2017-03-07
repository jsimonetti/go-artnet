package code

// SwRemote represents the trigger values ff the Node supports remote trigger inputs.
// The Node is responsible for ‘debouncing’ inputs. When the ArtPollReply is set to transmit
// automatically, (TalkToMe Bit 1), the ArtPollReply will be sent on both key down and key up
// events. However, the Controller should not assume that only one bit position has changed.
// The Remote inputs are used for remote event triggering or cueing
type SwRemote uint8

// WithRemote1 sets Remote 1 active
func (s SwRemote) WithRemote1(enable bool) SwRemote {
	if enable {
		return s | (1 << 0)
	}
	return s | (0 << 0)
}

// Remote1 indicates Remote 1 active
func (s SwRemote) Remote1() bool {
	return s&(1<<0) > 0
}

// WithRemote2 sets Remote 2 active
func (s SwRemote) WithRemote2(enable bool) SwRemote {
	if enable {
		return s | (1 << 1)
	}
	return s | (0 << 1)
}

// Remote2 indicates Remote 2 active
func (s SwRemote) Remote2() bool {
	return s&(1<<1) > 0
}

// WithRemote3 sets Remote 3 active
func (s SwRemote) WithRemote3(enable bool) SwRemote {
	if enable {
		return s | (1 << 2)
	}
	return s | (0 << 2)
}

// Remote3 indicates Remote 3 active
func (s SwRemote) Remote3() bool {
	return s&(1<<2) > 0
}

// WithRemote4 sets Remote 4 active
func (s SwRemote) WithRemote4(enable bool) SwRemote {
	if enable {
		return s | (1 << 3)
	}
	return s | (0 << 3)
}

// Remote4 indicates Remote 4 active
func (s SwRemote) Remote4() bool {
	return s&(1<<3) > 0
}

// WithRemote5 sets Remote 5 active
func (s SwRemote) WithRemote5(enable bool) SwRemote {
	if enable {
		return s | (1 << 4)
	}
	return s | (0 << 4)
}

// Remote5 indicates Remote 5 active
func (s SwRemote) Remote5() bool {
	return s&(1<<4) > 0
}

// WithRemote6 sets Remote 6 active
func (s SwRemote) WithRemote6(enable bool) SwRemote {
	if enable {
		return s | (1 << 5)
	}
	return s | (0 << 5)
}

// Remote6 indicates Remote 6 active
func (s SwRemote) Remote6() bool {
	return s&(1<<5) > 0
}

// WithRemote7 sets Remote 7 active
func (s SwRemote) WithRemote7(enable bool) SwRemote {
	if enable {
		return s | (1 << 6)
	}
	return s | (0 << 6)
}

// Remote7 indicates Remote 7 active
func (s SwRemote) Remote7() bool {
	return s&(1<<6) > 0
}

// WithRemote8 sets Remote 8 active
func (s SwRemote) WithRemote8(enable bool) SwRemote {
	if enable {
		return s | (1 << 7)
	}
	return s | (0 << 7)
}

// Remote8 indicates Remote 8 active
func (s SwRemote) Remote8() bool {
	return s&(1<<7) > 0
}

// String returns a string representation of SwRemote
func (s SwRemote) String() string {
	m1, m2, m3, m4, m5, m6, m7, m8 := "off", "off", "off", "off", "off", "off", "off", "off"
	if s.Remote1() {
		m1 = "on"
	}
	if s.Remote2() {
		m2 = "on"
	}
	if s.Remote3() {
		m3 = "on"
	}
	if s.Remote4() {
		m4 = "on"
	}
	if s.Remote5() {
		m5 = "on"
	}
	if s.Remote6() {
		m6 = "on"
	}
	if s.Remote7() {
		m7 = "on"
	}
	if s.Remote8() {
		m8 = "on"
	}

	return "SwRemote: Remote1: " + m1 + ", Remote2: " + m2 + ", Remote3: " + m3 + ", Remote4: " + m4 + ", Remote5: " + m5 + ", Remote6: " + m6 + ", Remote7: " + m7 + ", Remote8: " + m8
}
