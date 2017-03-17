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

	switch h.OpCode {
	case code.OpPoll:
		p = &ArtPollPacket{}
	case code.OpCode(swapUint16(uint16(code.OpPollReply))):
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
	case code.OpTodRequest:
	case code.OpTodData:
	case code.OpTodControl:
	case code.OpRdm:
	case code.OpRdmSub:
	case code.OpMedia:
	case code.OpMediaPatch:
	case code.OpMediaControl:
	case code.OpMediaContrlReply:
	case code.OpTimeCode:
		p = &ArtTimeCodePacket{}
	case code.OpTimeSync:
	case code.OpTrigger:
		p = &ArtTriggerPacket{}
	case code.OpDirectory:
	case code.OpDirectoryReply:
	case code.OpVideoSetup:
	case code.OpVideoPalette:
	case code.OpVideoData:
	case code.OpMacMaster:
	case code.OpMacSlave:
	case code.OpFirmwareMaster:
	case code.OpFirmwareReply:
	case code.OpFileTnMaster:
	case code.OpFileFnMaster:
	case code.OpFileFnReply:
	case code.OpIPProg:
		p = &ArtIPProgPacket{}
	case code.OpIPProgReply:
		p = &ArtIPProgReplyPacket{}
	default:
		return nil, fmt.Errorf("unimplemented opcode %#v found", h.OpCode)
	}

	err = p.UnmarshalBinary(b)
	return
}
