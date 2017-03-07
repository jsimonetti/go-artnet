package version

import (
	"testing"
)

func TestVersion(t *testing.T) {
	var tests = []struct {
		name string
		b    [2]byte
	}{
		{
			name: "V14",
			b:    [2]byte{0x0, 0x14},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bytes()
			if want, got := tt.b, b; want != got {
				t.Fatalf("unexpected Version bytes:\n- want: [%#v]\n-  got: [%#v]", want, got)
			}
		})
	}
}
