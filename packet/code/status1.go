package code

// Status1 is a general Status register containing bit fields.
type Status1 uint8

// WithUBEA sets the presence of UBEA
// false: UBEA not present or corrupt
// true:  UBEA present
func (s Status1) WithUBEA(enable bool) Status1 {
	if enable {
		return s | (1 << 0)
	}
	return s | (0 << 0)
}

// UBEA returns the status of the bit 0
func (s Status1) UBEA() bool {
	return s&(1<<0) > 0
}

// WithRDM sets the capability of the node for RDM
// false: Not capable of Remote Device Management (RDM).
// true:  Capable of Remote Device Management (RDM).
func (s Status1) WithRDM(enable bool) Status1 {
	if enable {
		return s | (1 << 1)
	}
	return s | (0 << 1)
}

// RDM returns the status of the bit 1
func (s Status1) RDM() bool {
	return s&(1<<1) > 0
}

// WithBootROM sets the boot location
// false: Normal firmware boot (from flash). Nodes that do not support dual boot,
//        clear this field to zero.
// true:  Booted from ROM.
func (s Status1) WithBootROM(enable bool) Status1 {
	if enable {
		return s | (1 << 2)
	}
	return s | (0 << 2)
}

// BootROM returns the status of the bit 2
func (s Status1) BootROM() bool {
	return s&(1<<2) > 0
}

// WithPortAddr sets the Port-Address Programming Authority
// v = "unknown": Port-Address Programming Authority unknown
// v = "front":   All Port-Address set by front panel controls.
// v = "net":     All or part of Port-Address programmed by networkor Web browser.
// v = "unused":  Not used.
func (s Status1) WithPortAddr(v string) Status1 {
	switch v {
	case "unknown":
		return s | (0 << 4)
	case "front":
		return s | (1 << 4)
	case "net":
		return s | (2 << 4)
	case "unused":
		return s | (3 << 4)
	}
	return s
}

// PortAddr returns the Port-Address Programming Authority
// "unknown": Port-Address Programming Authority unknown
// "front":   All Port-Address set by front panel controls.
// "net":     All or part of Port-Address programmed by networkor Web browser.
// "unused":  Not used.
func (s Status1) PortAddr() string {
	switch (s & (8 + 16)) >> 4 {
	case 0:
		return "unknown"
	case 1:
		return "front"
	case 2:
		return "net"
	case 3:
		return "unused"
	}
	return "error"
}

// WithIndicator sets Indicator state
// v = "unknown": Indicator state unknown.
// v = "locate":  Indicators in Locate / Identify Mode.
// v = "mute":    Indicators in Mute Mode.
// v = "normal":  Indicators in Normal Mode.
func (s Status1) WithIndicator(v string) Status1 {
	switch v {
	case "unknown":
		return s | (0 << 6)
	case "locate":
		return s | (1 << 6)
	case "mute":
		return s | (2 << 6)
	case "normal":
		return s | (3 << 6)
	}
	return s
}

// Indicator returns the Port-Address Programming Authority
// "unknown": Indicator state unknown.
// "locate":  Indicators in Locate / Identify Mode.
// "mute":    Indicators in Mute Mode.
// "normal":  Indicators in Normal Mode.
func (s Status1) Indicator() string {
	switch (s & (32 + 64)) >> 6 {
	case 0:
		return "unknown"
	case 1:
		return "locate"
	case 2:
		return "mutenet"
	case 3:
		return "normal"
	}
	return "error"
}

// String returns a string representation of Status1
func (s Status1) String() string {
	ubea, rdm, rom := "no", "no", "no"
	if s.UBEA() {
		ubea = "yes"
	}
	if s.RDM() {
		rdm = "yes"
	}
	if s.BootROM() {
		rom = "yes"
	}
	portAddr := s.PortAddr()
	indicator := s.Indicator()

	return "Status1: UBEA: " + ubea + ", RDM: " + rdm + ", BootRom: " + rom + ", PortAddr: " + portAddr + ", Indicator: " + indicator
}
