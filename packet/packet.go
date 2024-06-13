package packet

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"

	"github.com/jsimonetti/go-artnet/packet/code"
)

// Unmarshal will unmarshal the bytes into an ArtNetPacket
func Unmarshal(b []byte) (p ArtNetPacket, err error) {
	h := Header{}
	err = h.unmarshal(b)
	if err != nil {
		return
	}

	notImplErr := fmt.Errorf("unimplemented opcode %#v found", h.OpCode)

	switch h.GetOpCode() {
	case code.OpPoll:
		p = &ArtPollPacket{}
	case code.OpPollReply:
		p = &ArtPollReplyPacket{}
	case code.OpDiagData:
		p = &ArtDiagDataPacket{}
	case code.OpCommand:
		p = &ArtCommandPacket{}
	// OpOutput and OpDMX are the same, OpDMX is more common
	case code.OpDMX:
		p = &ArtDMXPacket{}
	case code.OpNzs:
		p = &ArtNzsPacket{}
	case code.OpSync:
		p = &ArtSyncPacket{}
	case code.OpAddress:
		p = &ArtAddressPacket{}
	case code.OpInput:
		return nil, notImplErr
	case code.OpTimeCode:
		p = &ArtTimeCodePacket{}
	case code.OpTrigger:
		p = &ArtTriggerPacket{}
	case code.OpIPProg:
		p = &ArtIPProgPacket{}
	case code.OpIPProgReply:
		p = &ArtIPProgReplyPacket{}
	case
		code.OpDirectory,
		code.OpDirectoryReply,
		code.OpFileFnMaster,
		code.OpFileFnReply,
		code.OpFileTnMaster,
		code.OpFirmwareMaster,
		code.OpFirmwareReply,
		code.OpMacMaster,
		code.OpMacSlave,
		code.OpMedia,
		code.OpMediaContrlReply,
		code.OpMediaControl,
		code.OpMediaPatch,
		code.OpRdm,
		code.OpRdmSub,
		code.OpTimeSync,
		code.OpTodControl,
		code.OpTodData,
		code.OpVideoData,
		code.OpVideoPalette,
		code.OpVideoSetup,
		code.OpTodRequest:
		return nil, fmt.Errorf("%w %#v", errNotImplementedOpCode, h.OpCode)
	default:
		return nil, fmt.Errorf("%w %#v", errInvalidOpCode, h.OpCode)
	}

	err = p.UnmarshalBinary(b)
	return
}

// ArtNetPacket is the interface used for passing around different kinds of ArtNet packets.
type ArtNetPacket interface {
	GetOpCode() code.OpCode

	encoding.BinaryMarshaler

	encoding.BinaryUnmarshaler
}

func marshalPacket(p ArtNetPacket) ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// checkPadAndUnmarshalPacket returns an error if len(b) is less than min, more than max
// and ensure len(b) is max by padding with bytes
func checkPadAndUnmarshalPacket(p ArtNetPacket, b []byte, min, max int) error {
	if len(b) < min {
		return errInvalidPacketMin
	}

	if len(b) > max {
		return errInvalidPacketMax
	}

	padding := make([]byte, max-len(b))
	b = append(b, padding...)
	return unmarshalPacket(p, b)

}

// unmarshalPacket fills p with from b (BigEndian)
// some packets will need to flip some values around after.
func unmarshalPacket(p ArtNetPacket, b []byte) error {
	buf := bytes.NewReader(b)
	return binary.Read(buf, binary.BigEndian, p)
}

func swapUint16(x uint16) uint16 {
	return x>>8 + x<<8
}
