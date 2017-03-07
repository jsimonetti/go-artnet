package code

// GoodInput indicates input status of the node
type GoodInput uint8

// WithReceive sets if Receive errors detected
func (s GoodInput) WithReceive(enable bool) GoodInput {
	if enable {
		return s | (1 << 2)
	}
	return s | (0 << 2)
}

// Receive indicates if Receive errors detected
func (s GoodInput) Receive() bool {
	return s&(1<<2) > 0
}

// WithDisabled sets Input is disabled
func (s GoodInput) WithDisabled(enable bool) GoodInput {
	if enable {
		return s | (1 << 3)
	}
	return s | (0 << 3)
}

// Disabled indicates Input is disabled
func (s GoodInput) Disabled() bool {
	return s&(1<<3) > 0
}

// WithText sets Channel includes DMX512 text packets
func (s GoodInput) WithText(enable bool) GoodInput {
	if enable {
		return s | (1 << 4)
	}
	return s | (0 << 4)
}

// Text indicates Channel includes DMX512 text packets
func (s GoodInput) Text() bool {
	return s&(1<<4) > 0
}

// WithSIP sets Channel includes DMX512 SIP’s
func (s GoodInput) WithSIP(enable bool) GoodInput {
	if enable {
		return s | (1 << 5)
	}
	return s | (0 << 5)
}

// SIP indicates Channel includes DMX512 SIP’s
func (s GoodInput) SIP() bool {
	return s&(1<<5) > 0
}

// WithTest sets Channel includes DMX512 test packets
func (s GoodInput) WithTest(enable bool) GoodInput {
	if enable {
		return s | (1 << 6)
	}
	return s | (0 << 6)
}

// Test indicates Channel includes DMX512 test packets
func (s GoodInput) Test() bool {
	return s&(1<<6) > 0
}

// WithData sets Data received
func (s GoodInput) WithData(enable bool) GoodInput {
	if enable {
		return s | (1 << 7)
	}
	return s | (0 << 7)
}

// Data indicates Data received
func (s GoodInput) Data() bool {
	return s&(1<<7) > 0
}

// String returns a string representation of GoodInput
func (s GoodInput) String() string {
	receive, disabled, text, sip, test, data := "no", "no", "no", "no", "no", "no"
	if s.Receive() {
		receive = "yes"
	}
	if s.Disabled() {
		disabled = "yes"
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

	return "GoodInput: ReceiveErrors: " + receive + ", Disabled: " + disabled + ", DMX512Text: " + text + ", DMX512SIP: " + sip + ", DMX512Test: " + test + ", DataReceived: " + data
}
