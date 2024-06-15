package artnettypes

// Version 14 support for this package only
// The hex representation of this has the high byte first
// because of the representation of Version
const V14 Version = 0x000e

const CurrentVersion Version = V14

type Version uint16

// GetVersion returns the version parsed by validate method
func (v Version) GetVersion() Version {
	return Version(swapUint16(uint16(v)))
}

func swapUint16(x uint16) uint16 {
	return x>>8 + x<<8
}
