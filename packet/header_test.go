package packet

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"

	"github.com/jsimonetti/go-artnet/packet/code"
)

func TestHeaderValidate(t *testing.T) {
	v := version
	v1 := byte(v & 0x00ff)
	v0 := byte((v & 0xff00) >> 8)

	cases := []struct {
		pkg []byte
		op  code.OpCode
		err string
	}{
		{pkg: makePkg(t, code.OpPoll, v0, v1), op: code.OpPoll, err: ""},
		{pkg: makePkg(t, code.OpPoll, v0, 0x00), op: code.OpPoll, err: "incompatible version. want: 14, got: 0"},
		{pkg: makePkg(t, code.OpPoll, v0, 0x0f), op: code.OpPoll, err: ""},
	}

	for i, c := range cases {
		h := Header{}
		err := h.unmarshal(c.pkg)
		if err == nil {
			err = h.validate(c.op)
		}
		if (err == nil && c.err != "") || (c.err == "" && err != nil) || (c.err != "" && err != nil && c.err != err.Error()) {
			t.Errorf("case %d: Expected to get err %q, got %q", i, c.err, err)
		}
	}
}

func makePkg(t *testing.T, opCode code.OpCode, data ...byte) (pkg []byte) {
	pkg = append(pkg, artNet[0], artNet[1], artNet[2], artNet[3], artNet[4], artNet[5], artNet[6], artNet[7])
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

func TestHeaderWithoutVersionValidate(t *testing.T) {

	cases := []struct {
		pkg []byte
		err string
	}{
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
		h := &HeaderWithoutVersion{}
		buf := bytes.NewReader(c.pkg)
		err := binary.Read(buf, binary.BigEndian, h)
		if err == nil {
			err = h.validate(code.OpPollReply)
		}
		if (err == nil && c.err != "") || (c.err == "" && err != nil) || (c.err != "" && err != nil && c.err != err.Error()) {
			t.Errorf("case %d: Expected to get err %q, got %q", i, c.err, err)
		}
	}
}
