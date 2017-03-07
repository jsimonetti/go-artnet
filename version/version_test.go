package version

import (
	"bytes"
	"testing"
)

func TestVersion(t *testing.T) {
	var tests = []struct {
		name string
		b    []byte
	}{
		{
			name: "V14",
			b:    []byte{0x0, 0x14},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bytes()
			if want, got := tt.b, b; !bytes.Equal(want, got) {
				t.Fatalf("unexpected Version bytes:\n- want: [%#v]\n-  got: [%#v]", want, got)
			}
		})
	}
}
