package code

// PortType defines the operation and protocol of each channel. (A product with 4 inputs
// and 4 outputs would report 0xc0, 0xc0, 0xc0, 0xc0). The array length is fixed, independent
// of the number of inputs or outputs physically available on the Node
type PortType uint8

// WithInput sets channel can input from the Art-Net Network
func (s PortType) WithInput(enable bool) PortType {
	if enable {
		return s | (1 << 6)
	}
	return s | (0 << 6)
}

// Input indicates channel can input from the Art-Net Network
func (s PortType) Input() bool {
	return s&(1<<6) > 0
}

// WithOutput sets channel can output onto the Art-Net Network
func (s PortType) WithOutput(enable bool) PortType {
	if enable {
		return s | (1 << 7)
	}
	return s | (0 << 7)
}

// Output indicates channel can output onto the Art-Net Network
func (s PortType) Output() bool {
	return s&(1<<7) > 0
}

// WithType sets the Port-Address Programming Authority
// v = "DMX512" | "MIDI" | "Avab" | "Colortran CMX" | "ADB 62.5" | "Art-Net"
func (s PortType) WithType(v string) PortType {
	switch v {
	case "DMX512":
		return s | (0 << 0) | (1 << 0) | (2 << 0) | (3 << 0) | (4 << 0) | (5 << 0)
	case "MIDI":
		return s | (0 << 0) | (1 << 0) | (2 << 0) | (3 << 0) | (4 << 0) | (5 << 1)
	case "Avab":
		return s | (0 << 0) | (1 << 0) | (2 << 0) | (3 << 0) | (4 << 1) | (5 << 0)
	case "Colortran CMX":
		return s | (0 << 0) | (1 << 0) | (2 << 0) | (3 << 0) | (4 << 1) | (5 << 1)
	case "ADB 62.5":
		return s | (0 << 0) | (1 << 0) | (2 << 0) | (3 << 1) | (4 << 0) | (5 << 0)
	case "Art-Net":
		return s | (0 << 0) | (1 << 0) | (2 << 0) | (3 << 1) | (4 << 0) | (5 << 1)
	}
	return s
}

// Type returns the Port-Address Programming Authority
// "DMX512" | "MIDI" | "Avab" | "Colortran CMX" | "ADB 62.5" | "Art-Net"
func (s PortType) Type() string {
	switch s >> 0 {
	case 0:
		return "DMX512"
	case 1:
		return "MIDI"
	case 2:
		return "Avab"
	case 3:
		return "Colortran CMX"
	case 4:
		return "ADB 62.5"
	case 5:
		return "Art-Net"
	}
	return "error"
}

// String returns a string representation of PortType
func (s PortType) String() string {
	output, input := "no", "no"
	if s.Output() {
		output = "yes"
	}
	if s.Input() {
		input = "yes"
	}

	porttype := s.Type()
	return "PortType: Type: " + porttype + ", Output: " + output + ", Input: " + input
}
