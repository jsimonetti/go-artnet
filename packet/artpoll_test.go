package packet

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

func TestArtPollPacketMarshal(t *testing.T) {
	tests := []struct {
		name string
		p    ArtPollPacket
		b    []byte
		err  error
	}{
		{
			name: "Empty",
			p: ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpPoll,
					Version: version.Bytes(),
				},
			},
			b: []byte{
				0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
				0x00, 0x20, 0x00, 0x0e, 0x00, 0x00,
			},
		},
		{
			name: "WithDiagnosticsPrioLow",
			p: ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpPoll,
					Version: version.Bytes(),
				},
				TalkToMe: new(code.TalkToMe).WithDiagnostics(true),
				Priority: code.DpLow,
			},
			b: []byte{
				0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
				0x00, 0x20, 0x00, 0x0e, 0x004, 0x10,
			},
		},
		{
			name: "WithDiagnosticsUniPrioMedium",
			p: ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpPoll,
					Version: version.Bytes(),
				},
				TalkToMe: new(code.TalkToMe).WithDiagnostics(true).WithDiagUnicast(true),
				Priority: code.DpMed,
			},
			b: []byte{
				0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
				0x00, 0x20, 0x00, 0x0e, 0x00c, 0x40,
			},
		},
		{
			name: "WithReplyOnChangeVlcPrioVolatile",
			p: ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpPoll,
					Version: version.Bytes(),
				},
				TalkToMe: new(code.TalkToMe).WithReplyOnChange(true).WithVLC(true),
				Priority: code.DpVolatile,
			},
			b: []byte{
				0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
				0x00, 0x20, 0x00, 0x0e, 0x12, 0xf0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.p.MarshalBinary()

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
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

func TestArtPollPacketUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		p    ArtPollPacket
		b    []byte
		err  error
	}{
		{
			name: "Empty",
			p: ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpPoll,
					Version: version.Bytes(),
				},
			},
			b: []byte{
				0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
				0x00, 0x20, 0x00, 0x0e, 0x00, 0x00,
			},
		},
		{
			name: "WithDiagnosticsPrioLow",
			p: ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpPoll,
					Version: version.Bytes(),
				},
				TalkToMe: new(code.TalkToMe).WithDiagnostics(true),
				Priority: code.DpLow,
			},
			b: []byte{
				0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
				0x00, 0x20, 0x00, 0x0e, 0x004, 0x10,
			},
		},
		{
			name: "WithDiagnosticsUniPrioMedium",
			p: ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpPoll,
					Version: version.Bytes(),
				},
				TalkToMe: new(code.TalkToMe).WithDiagnostics(true).WithDiagUnicast(true),
				Priority: code.DpMed,
			},
			b: []byte{
				0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
				0x00, 0x20, 0x00, 0x0e, 0x00c, 0x40,
			},
		},
		{
			name: "WithReplyOnChangeVlcPrioVolatile",
			p: ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpPoll,
					Version: version.Bytes(),
				},
				TalkToMe: new(code.TalkToMe).WithReplyOnChange(true).WithVLC(true),
				Priority: code.DpVolatile,
			},
			b: []byte{
				0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
				0x00, 0x20, 0x00, 0x0e, 0x12, 0xf0,
			},
		},
		{
			name: "LiveCapture",
			p: ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpPoll,
					Version: [2]uint8{0x0, 0xe},
				},
				TalkToMe: new(code.TalkToMe).WithReplyOnChange(false),
				Priority: code.DpAll,
			},
			b: []byte{
				0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00,
				0x00, 0x20, 0x00, 0x0e, 0x00, 0x00,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a ArtPollPacket
			err := a.UnmarshalBinary(tt.b)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
			}
			if err != nil {
				return
			}

			if want, got := tt.p, a; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Message bytes:\n- want: [%#v]\n-  got: [%#v]", want, got)
			}
		})
	}
}
