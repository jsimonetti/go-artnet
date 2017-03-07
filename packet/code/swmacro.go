package code

// SwMacro represents the trigger values if the Node supports macro key inputs. The Node is
// responsible for ‘debouncing’ inputs. When the ArtPollReply is set to transmit automatically,
// (TalkToMe Bit 1), the ArtPollReply will be sent on both key down and key up events. However, the
// Controller should not assume that only one bit position has changed. The Macro inputs are used
// for remote event triggering or cueing.
type SwMacro uint8

// WithMacro1 sets Macro 1 active
func (s SwMacro) WithMacro1(enable bool) SwMacro {
	if enable {
		return s | (1 << 0)
	}
	return s | (0 << 0)
}

// Macro1 indicates Macro 1 active
func (s SwMacro) Macro1() bool {
	return s&(1<<0) > 0
}

// WithMacro2 sets Macro 2 active
func (s SwMacro) WithMacro2(enable bool) SwMacro {
	if enable {
		return s | (1 << 1)
	}
	return s | (0 << 1)
}

// Macro2 indicates Macro 2 active
func (s SwMacro) Macro2() bool {
	return s&(1<<1) > 0
}

// WithMacro3 sets Macro 3 active
func (s SwMacro) WithMacro3(enable bool) SwMacro {
	if enable {
		return s | (1 << 2)
	}
	return s | (0 << 2)
}

// Macro3 indicates Macro 3 active
func (s SwMacro) Macro3() bool {
	return s&(1<<2) > 0
}

// WithMacro4 sets Macro 4 active
func (s SwMacro) WithMacro4(enable bool) SwMacro {
	if enable {
		return s | (1 << 3)
	}
	return s | (0 << 3)
}

// Macro4 indicates Macro 4 active
func (s SwMacro) Macro4() bool {
	return s&(1<<3) > 0
}

// WithMacro5 sets Macro 5 active
func (s SwMacro) WithMacro5(enable bool) SwMacro {
	if enable {
		return s | (1 << 4)
	}
	return s | (0 << 4)
}

// Macro5 indicates Macro 5 active
func (s SwMacro) Macro5() bool {
	return s&(1<<4) > 0
}

// WithMacro6 sets Macro 6 active
func (s SwMacro) WithMacro6(enable bool) SwMacro {
	if enable {
		return s | (1 << 5)
	}
	return s | (0 << 5)
}

// Macro6 indicates Macro 6 active
func (s SwMacro) Macro6() bool {
	return s&(1<<5) > 0
}

// WithMacro7 sets Macro 7 active
func (s SwMacro) WithMacro7(enable bool) SwMacro {
	if enable {
		return s | (1 << 6)
	}
	return s | (0 << 6)
}

// Macro7 indicates Macro 7 active
func (s SwMacro) Macro7() bool {
	return s&(1<<6) > 0
}

// WithMacro8 sets Macro 8 active
func (s SwMacro) WithMacro8(enable bool) SwMacro {
	if enable {
		return s | (1 << 7)
	}
	return s | (0 << 7)
}

// Macro8 indicates Macro 8 active
func (s SwMacro) Macro8() bool {
	return s&(1<<7) > 0
}

// String returns a string representation of SwMacro
func (s SwMacro) String() string {
	m1, m2, m3, m4, m5, m6, m7, m8 := "off", "off", "off", "off", "off", "off", "off", "off"
	if s.Macro1() {
		m1 = "on"
	}
	if s.Macro2() {
		m2 = "on"
	}
	if s.Macro3() {
		m3 = "on"
	}
	if s.Macro4() {
		m4 = "on"
	}
	if s.Macro5() {
		m5 = "on"
	}
	if s.Macro6() {
		m6 = "on"
	}
	if s.Macro7() {
		m7 = "on"
	}
	if s.Macro8() {
		m8 = "on"
	}

	return "SwMacro: Macro1: " + m1 + ", Macro2: " + m2 + ", Macro3: " + m3 + ", Macro4: " + m4 + ", Macro5: " + m5 + ", Macro6: " + m6 + ", Macro7: " + m7 + ", Macro8: " + m8
}
