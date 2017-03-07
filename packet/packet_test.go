package packet

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/jsimonetti/artnet/packet/code"
)

var artNetPackets = []struct {
	name string
	p    Packet
	b    []byte
	me   error
	ue   error
}{
	{
		name: "Empty",
		p:    Packet{},
		b:    []byte{},
		ue:   errIncorrectHeaderLength,
		me:   errInvalidOpCode,
	},
	{
		name: "NotArt-Net",
		p:    Packet{},
		b: []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00,
		},
		ue: errInvalidPacket,
		me: errInvalidOpCode,
	},
	{
		name: "EmptyOpCode",
		p:    Packet{},
		b: []byte{
			0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
			0x00, 0x00,
		},
		ue: errInvalidOpCode,
		me: errInvalidOpCode,
	},
	{
		name: "WithOpCodeDMX",
		p: Packet{
			OpCode: code.OpDMX,
		},
		b: []byte{
			0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
			0x00, 0x50,
		},
	},
	{
		name: "WithOpCodeVideo",
		p: Packet{
			OpCode: code.OpVideoData,
		},
		b: []byte{
			0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
			0x40, 0xa0,
		},
	},
	{
		name: "WithTrailingData",
		p: Packet{
			OpCode: code.OpDMX,
			data:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		b: []byte{
			0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
			0x00, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
}

func TestPacketMarshal(t *testing.T) {
	for _, tt := range artNetPackets {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.p.MarshalBinary()
			if err != tt.me {
				t.Fatalf("unexpected error: want: %v, got: %v", tt.me, err)
			}

			if err != nil {
				return
			}

			if want, got := tt.b, b; !bytes.Equal(want, got) {
				t.Fatalf("unexpected Message bytes:\n- want: [%# x]\n-  got: [%# x]", want, got)
			}
		})
	}
}

func TestPacketUnmarshal(t *testing.T) {
	for _, tt := range artNetPackets {
		t.Run(tt.name, func(t *testing.T) {
			var p Packet
			err := p.UnmarshalBinary(tt.b)
			if err != tt.ue {
				t.Fatalf("unexpected error: want: %v, got: %v", tt.ue, err)
			}

			if err != nil {
				return
			}

			if want, got := tt.p, p; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Message bytes:\n- want: [%#v]\n-  got: [%#v]", want, got)
			}
		})
	}
}
