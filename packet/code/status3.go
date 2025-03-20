package code

// Status3 indicates general product state
type Status3 uint8

// WithFailSafe sets how the node behaves in the
// event that network data is lost.
// v = "last":  Hold last state.
// v = "zero":  All outputs to zero.
// v = "full":  All outputs to full.
// v = "scene": Playback failsafe scene.
func (s Status3) WithFailSafe(v string) Status3 {
	switch v {
	case "last":
		return s | (0 << 6)
	case "zero":
		return s | (1 << 6)
	case "full":
		return s | (2 << 6)
	case "scene":
		return s | (3 << 6)
	}
	return s
}

// WithFailSafe returns how the node will behave in the
// event that network data is lost.
// v = "last":  Hold last state.
// v = "zero":  All outputs to zero.
// v = "full":  All outputs to full.
// v = "scene": Playback failsafe scene.
func (s Status3) FailSafe() string {
	switch (s & (0xc0)) >> 6 {
	case 0:
		return "last"
	case 1:
		return "zero"
	case 2:
		return "full"
	case 3:
		return "scene"
	}
	return "error"
}

// WithRDMNet sets if product supports RDMNet
func (s Status3) WithRDMNet(enable bool) Status3 {
	if enable {
		return s | (1 << 2)
	}
	return s | (0 << 2)
}

// RDMNet indicates if product supports RDMNet
func (s Status3) RDMNet() bool {
	return s&(1<<2) > 0
}

// WithPortSwitching sets if product supports switching input and output ports
func (s Status3) WithPortSwitching(enable bool) Status3 {
	if enable {
		return s | (1 << 3)
	}
	return s | (0 << 3)
}

// PortSwitching indicates if product supports switching input and output ports
func (s Status3) PortSwitching() bool {
	return s&(1<<3) > 0
}

// WithLLRP sets if product supports LLRP
func (s Status3) WithLLRP(enable bool) Status3 {
	if enable {
		return s | (1 << 4)
	}
	return s | (0 << 4)
}

// LLRP indicates if product supports LLRP
func (s Status3) LLRP() bool {
	return s&(1<<4) > 0
}

// WithFailOver sets if product supports fail-over
func (s Status3) WithFailOver(enable bool) Status3 {
	if enable {
		return s | (1 << 5)
	}
	return s | (0 << 5)
}

// FailOver indicates if product supports fail-over
func (s Status3) FailOver() bool {
	return s&(1<<5) > 0
}

// String returns a string representation of Status3
func (s Status3) String() string {
	rdmNet, portSwitching, llrp, failOver := "no", "no", "no", "no"
	if s.RDMNet() {
		rdmNet = "yes"
	}
	if s.PortSwitching() {
		portSwitching = "yes"
	}
	if s.LLRP() {
		llrp = "yes"
	}
	if s.FailOver() {
		failOver = "yes"
	}
	failSafe := s.FailSafe()

	return "Status3: FailSafe: " + failSafe + "RDMNet: " + rdmNet + ", Port Switching: " + portSwitching + ", LLRP: " + llrp + ", FailOvering: " + failOver
}
