package packet

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"

	"github.com/jsimonetti/go-artnet/packet/code"
	"github.com/jsimonetti/go-artnet/version"
)

func TestHeaderValidate(t *testing.T) {
	v := version.Bytes()

	cases := []struct {
		pkg []byte
		err string
	}{
		{pkg: makePkg(t, code.OpPoll, v[0], v[1]), err: ""},
		{pkg: makePkg(t, code.OpPoll, v[0], 0x00), err: "incompatible version. want: =>14, got: 0"},
		{pkg: makePkg(t, code.OpPoll, v[0], 0x0f), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "2.0.0.10")...), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "2.14.0.10")...), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "2.16.0.10")...), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "2.100.0.10")...), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "192.168.0.10")...), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "10.0.0.10")...), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "10.10.0.10")...), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "10.20.0.10")...), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "10.30.0.10")...), err: ""},
		{pkg: makePkg(t, code.OpPollReply, ipToBytes(t, "10.30.0.10")...), err: ""},
	}

	for i, c := range cases {
		h := Header{}
		err := h.unmarshal(c.pkg)
		if (err == nil && c.err != "") || (c.err == "" && err != nil) || (c.err != "" && err != nil && c.err != err.Error()) {
			t.Errorf("case %d: Expected to get err %q, got %q", i, c.err, err)
		}
	}
}

func makePkg(t *testing.T, opCode code.OpCode, data ...byte) (pkg []byte) {
	pkg = append(pkg, ArtNet[0], ArtNet[1], ArtNet[2], ArtNet[3], ArtNet[4], ArtNet[5], ArtNet[6], ArtNet[7])
	pkg = append(pkg, opCodeToBytes(t, opCode)...)
	pkg = append(pkg, data...)

	return
}

func opCodeToBytes(t *testing.T, opCode code.OpCode) []byte {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.LittleEndian, opCode); err != nil {
		t.Fatalf("failed to encode opCode: %v", err)
	}
	return buf.Bytes()
}

func ipToBytes(t *testing.T, ipStr string) []byte {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		t.Fatalf("%q is not an IP, thanks", ipStr)
	}
	return []byte(ip)[12:]
}
