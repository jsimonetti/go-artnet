package code

// GoodOutput indicates output status of the node
type GoodOutput uint8

// WithACN sets Output to transmit sACN
func (s GoodOutput) WithACN(enable bool) GoodOutput {
	if enable {
		return s | (1 << 0)
	}
	return s | (0 << 0)
}

// ACN indicates Output is selected to transmit sACN
func (s GoodOutput) ACN() bool {
	return s&(1<<0) > 0
}

// WithLTP sets Merge Mode is LTP
func (s GoodOutput) WithLTP(enable bool) GoodOutput {
	if enable {
		return s | (1 << 1)
	}
	return s | (0 << 1)
}

// LTP indicates Merge Mode is LTP
func (s GoodOutput) LTP() bool {
	return s&(1<<1) > 0
}

// WithOutput sets DMX output short detected on power up
func (s GoodOutput) WithOutput(enable bool) GoodOutput {
	if enable {
		return s | (1 << 2)
	}
	return s | (0 << 2)
}

// Output indicates DMX output short detected on power up
func (s GoodOutput) Output() bool {
	return s&(1<<2) > 0
}

// WithMerging sets Output is merging ArtNet data
func (s GoodOutput) WithMerging(enable bool) GoodOutput {
	if enable {
		return s | (1 << 3)
	}
	return s | (0 << 3)
}

// Merging indicates Output is merging ArtNet data
func (s GoodOutput) Merging() bool {
	return s&(1<<3) > 0
}

// WithText sets Channel includes DMX512 text packets
func (s GoodOutput) WithText(enable bool) GoodOutput {
	if enable {
		return s | (1 << 4)
	}
	return s | (0 << 4)
}

// Text indicates Channel includes DMX512 text packets
func (s GoodOutput) Text() bool {
	return s&(1<<4) > 0
}

// WithSIP sets Channel includes DMX512 SIP’s
func (s GoodOutput) WithSIP(enable bool) GoodOutput {
	if enable {
		return s | (1 << 5)
	}
	return s | (0 << 5)
}

// SIP indicates Channel includes DMX512 SIP’s
func (s GoodOutput) SIP() bool {
	return s&(1<<5) > 0
}

// WithTest sets Channel includes DMX512 test packets
func (s GoodOutput) WithTest(enable bool) GoodOutput {
	if enable {
		return s | (1 << 6)
	}
	return s | (0 << 6)
}

// Test indicates Channel includes DMX512 test packets
func (s GoodOutput) Test() bool {
	return s&(1<<6) > 0
}

// WithData sets Data transmitted
func (s GoodOutput) WithData(enable bool) GoodOutput {
	if enable {
		return s | (1 << 7)
	}
	return s | (0 << 7)
}

// Data indicates Data transmitted
func (s GoodOutput) Data() bool {
	return s&(1<<7) > 0
}

// String returns a string representation of TalkToMe
func (s GoodOutput) String() string {
	acn, ltp, output, merging, text, sip, test, data := "no", "no", "no", "no", "no", "no", "no", "no"
	if s.LTP() {
		ltp = "yes"
	}
	if s.ACN() {
		acn = "yes"
	}
	if s.Output() {
		output = "yes"
	}
	if s.Merging() {
		merging = "yes"
	}
	if s.Text() {
		text = "yes"
	}
	if s.SIP() {
		sip = "yes"
	}
	if s.Test() {
		test = "yes"
	}
	if s.Data() {
		data = "yes"
	}

	return "GoodInput: OutputACN: " + acn + ", LTPMergeMode: " + ltp + ", OutputShort: " + output + ", Merging: " + merging + ", DMX512Text: " + text + ", DMX512SIP: " + sip + ", DMX512Test: " + test + ", DataReceived: " + data
}
