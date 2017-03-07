package code

//go:generate stringer -type=StyleCode

// StyleCode defines the class of data following an UDP packet.
type StyleCode uint8

// Valid returns wether the OpCode is valid
func ValidStyle(o StyleCode) bool {
	if o >= StyleCode(len(_StyleCode_index)-1) {
		return false
	}
	return true
}

const (
	// StNode A DMX to / from Art-Net device
	StNode StyleCode = 0x00

	// StController A lighting console.
	StController StyleCode = 0x01

	// StMedia A Media Server.
	StMedia StyleCode = 0x02

	// StRoute A network routing device.
	StRoute StyleCode = 0x03

	// StBackup A backup device.
	StBackup StyleCode = 0x04

	// StConfig A configuration or diagnostic tool.
	StConfig StyleCode = 0x05

	// StVisual A visualiser.
	StVisual StyleCode = 0x06
)
