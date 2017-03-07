package code

// TalkToMe sets the behaviour of a Node
// only bits 1-4 matter, rest is zero
type TalkToMe uint8

// WithReplyOnChange allows the Controller to be informed of changes
// without the need to continuously poll.
func (t TalkToMe) WithReplyOnChange(enable bool) TalkToMe {
	if enable {
		return t | (1 << 1)
	}
	return t | (0 << 1)
}

// ReplyOnChange returns the status of the bit 1
func (t TalkToMe) ReplyOnChange() bool {
	return t&(1<<1) > 0
}

// WithDiagnostics sends diagnostics messages
func (t TalkToMe) WithDiagnostics(enable bool) TalkToMe {
	if enable {
		return t | (1 << 2)
	}
	return t | (0 << 2)
}

// Diagnostics returns the stats of the bit 2
func (t TalkToMe) Diagnostics() bool {
	return t&(1<<2) > 0
}

// WithDiagUnicast determines wether diagnostics messages are unicast or broadcast
func (t TalkToMe) WithDiagUnicast(enable bool) TalkToMe {
	if enable {
		return t | (1 << 3)
	}
	return t | (0 << 3)
}

// DiagUnicast returns the stats of the bit 3
func (t TalkToMe) DiagUnicast() bool {
	return t&(1<<3) > 0
}

// WithVLC enable or disable VLC transmission
func (t TalkToMe) WithVLC(enable bool) TalkToMe {
	if enable {
		return t | (1 << 4)
	}
	return t | (0 << 4)
}

// VLC returns the stats of the bit 4
func (t TalkToMe) VLC() bool {
	return t&(1<<4) > 0
}

// String returns a string representation of TalkToMe
func (t TalkToMe) String() string {
	roc, diag, uni, vlc := "no", "no", "no", "no"
	if t.ReplyOnChange() {
		roc = "yes"
	}
	if t.Diagnostics() {
		diag = "yes"
	}
	if t.DiagUnicast() {
		uni = "yes"
	}
	if t.VLC() {
		vlc = "yes"
	}

	return "TalkToMe: ReplyOnChange: " + roc + ", Diagnostics: " + diag + " (as unicast: " + uni + "), VLC: " + vlc
}
