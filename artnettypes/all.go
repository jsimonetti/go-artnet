package artnettypes

import "fmt"

type DMXData [512]byte

type BindIndex uint8

// Address contains a universe address
type Address struct {
	SubUni uint8
	Net    uint8 // 0-128
}

// String returns a string representation of Address
func (a Address) String() string {
	return fmt.Sprintf("%d:%d.%d", a.Net, a.SubUni>>4, a.SubUni&0x0f)
}

// Integer returns the integer representation of Address
func (a Address) Integer() int {
	return int(uint16(a.Net)<<8 | uint16(a.SubUni))
}

type ID [8]byte

// ArtNet is the fixed string "Art-Net" terminated with a zero
var ArtNet = ID{0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00}
