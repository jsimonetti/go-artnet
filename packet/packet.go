package packet

import (
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

	switch h.OpCode {
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
