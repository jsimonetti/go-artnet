package packet

import (
	"reflect"
	"testing"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

func TestUnmarshal(t *testing.T) {
	cases := []struct {
		b   []byte
		err error
		pkg ArtNetPacket
	}{
		{
			b: []byte{65, 114, 116, 45, 78, 101, 116, 0, 0, 32, 0, 14, 2, 0},
			pkg: &ArtPollPacket{
				Header: Header{
					ID:      ArtNet,
					Version: version.Bytes(),
					OpCode:  code.OpPoll,
				},
				TalkToMe: new(code.TalkToMe).WithReplyOnChange(true),
			},
		},
		{
			b: []byte{
				65, 114, 116, 45, 78, 101, 116, 0, 0, 33, 2, 0, 0, 20, 54, 25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				98, 97, 114, 97, 100, 100, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			pkg: &ArtPollReplyPacket{
				ID:               ArtNet,
				OpCode:           code.OpPollReply,
				IPAddress:        [4]uint8{2, 0, 0, 20},
				Port:             ArtNetPort,
				ESTAmanufacturer: [2]uint8{},
				ShortName:        [18]uint8{0x62, 0x61, 0x72, 0x61, 0x64, 0x64, 0x75, 0x72},
				LongName:         [64]uint8{},
				NodeReport:       code.NodeReport{},
				PortTypes:        [4]code.PortType{},
				GoodInput:        [4]code.GoodInput{},
				GoodOutput:       [4]code.GoodOutput{},
				SwIn:             [4]uint8{},
				SwOut:            [4]uint8{},
				Style:            code.StController,
				Macaddress:       [6]uint8{},
				BindIP:           [4]uint8{},
			},
		},
		{
			b: []byte{
				65, 114, 116, 45, 78, 101, 116, 0, 0, 80, 0, 14, 179, 0, 7, 0, 2, 0, 0, 0, 0, 20, 0, 0, 0, 20, 0,
				0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0,
				20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 20, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			pkg: &ArtDMXPacket{
				Header: Header{
					ID:      ArtNet,
					OpCode:  code.OpDMX,
					Version: version.Bytes(),
				},
				Sequence: 0xb3,
				Physical: 0x0,
				SubUni:   0x7,
				Net:      0x0,
				Length:   0x200,
				Data: [512]uint8{
					0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0,
					0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0,
					0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0,
					0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				}},
		},
	}

	for i, c := range cases {
		p, err := Unmarshal(c.b)
		if err != nil && err != c.err {
			t.Errorf("case %d: expected to get error %v, got %v", i, c.err, err)
			continue
		}

		if err == nil && c.err != nil {
			t.Errorf("case %d: expected to get error %v, got nil", i, err)
		}

		if !reflect.DeepEqual(c.pkg, p) {
			t.Errorf("case %d: expected to get packet %#v, got %#v", i, c.pkg, p)
		}
	}
}
