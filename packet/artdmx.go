package packet

import (
	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/types"
)

var _ ArtNetPacket = &ArtDMXPacket{}

// ArtDMXPacket contains an ArtDMX Packet.
//
// ArtDmx is the data packet used to transfer DMX512 data. The format is identical for
// Node to Controller, Node to Node and Controller to Node.
// The Data is output through the DMX O/P port corresponding to the Universe setting. In
// the absence of received ArtDmx packets, each DMX O/P port re-transmits the same
// frame continuously.
// The first complete DMX frame received at each input port is placed in an ArtDmx packet
// as above and transmitted as an ArtDmx packet containing the relevant Universe
// parameter. Each subsequent DMX frame containing new data (different length or
// different contents) is also transmitted as an ArtDmx packet.
// Nodes do not transmit ArtDmx for DMX512 inputs that have not received data since
// power on.
// However, an input that is active but not changing, will re-transmit the last valid ArtDmx
// packet at approximately 4-second intervals. (Note. In order to converge the needs of Art-Net
// and sACN it is recommended that Art-Net devices actually use a re-transmit time of
// 800mS to 1000mS).
// A DMX input that fails will not continue to transmit ArtDmx data.
//
// Packet Strategy:
//
//	Controller -  Receive:            Application Specific
//	              Unicast Transmit:   Yes
//	              Broadcast Transmit: No
//	Node -        Receive:            Application Specific
//	              Unicast Transmit:   Yes
//	              Broadcast Transmit: No
//	MediaServer - Receive:            Application Specific
//	              Unicast Transmit:   Yes
//	              Broadcast Transmit: No
type ArtDMXPacket struct {
	// Inherit the Header
	Header

	// Sequence number is used to ensure that ArtDmx packets are used in the correct order.
	// When Art-Net is carried over a medium such as the Internet, it is possible that ArtDmx packets
	// will reach the receiver out of order.
	// This field is incremented in the range 0x01 to 0xff to allow the receiving node to resequence
	// packets. The Sequence field is set to 0x00 to disable this feature
	Sequence uint8

	// Physical input port from which DMX512 data was input. This field is for information
	// only. Use Universe for data routing
	Physical uint8

	// SubUni is the low byte of the 15 bit Port-Address to which this packet is destined
	SubUni uint8

	// Net is the top 7 bits of the 15 bit Port-Address to which this packet is destined
	Net uint8

	// Length indicates the length of the data. This value should be an even number in the
	// range 2 – 512. It represents the number of DMX512 channels encoded in packet.
	// NB: Products which convert Art-Net to DMX512 may opt to always send 512 channels
	Length uint16

	// Data is a string of DMX512 lighting data
	Data types.DMXData
}

const (
	minimumArtDMXPacketSize int = 20
	maximumArtDMXPacketSize int = minimumArtDMXPacketSize + 510
)

// NewArtDMXPacket returns an ArtNetPacket with the correct OpCode
func NewArtDMXPacket(a types.Address, data types.DMXData, sequence uint8) *ArtDMXPacket {
	p := &ArtDMXPacket{
		Header:   NewHeader(code.OpDMX),
		Net:      a.Net,
		SubUni:   a.SubUni,
		Length:   512,
		Sequence: sequence,
	}

	copy(p.Data[:], data[:])

	return p
}

// MarshalBinary marshals an ArtDMXPacket into a byte slice.
func (p *ArtDMXPacket) MarshalBinary() ([]byte, error) {
	return marshalPacket(p)
}

// UnmarshalBinary unmarshals the contents of a byte slice into an ArtDMXPacket.
func (p *ArtDMXPacket) UnmarshalBinary(b []byte) error {
	if len(b)%2 != 0 {
		return errInvalidPacketBoundary
	}

	err := checkPadAndUnmarshalPacket(p, b, minimumArtDMXPacketSize, maximumArtDMXPacketSize)
	if err != nil {
		return err
	}

	return p.Header.validate(code.OpDMX)
}

func (p ArtDMXPacket) GetLength() uint16 {
	return swapUint16(p.Length)
}

func (p ArtDMXPacket) GetAddress() types.Address {
	return types.Address{
		Net:    p.Net,
		SubUni: p.SubUni,
	}
}
