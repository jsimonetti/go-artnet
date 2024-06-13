package types

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestVersion(t *testing.T) {
	var tests = []struct {
		name string
		b    [2]byte
	}{
		{
			name: "V14",
			b:    [2]byte{0x0, 0x0e},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := binary.Write(&buf, binary.BigEndian, CurrentVersion); err != nil {
				t.Fatalf("unexpected error: %e", err)
			}
			if want, got := tt.b, buf.Bytes(); want[0] != got[0] || want[1] != got[1] {
				t.Fatalf("unexpected Version bytes:\n- want: [%#v]\n-  got: [%#v]", want, got)
			}
		})
	}
}
