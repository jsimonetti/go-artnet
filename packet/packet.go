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
	case code.OpTodRequest:
		return nil, notImplErr
	case code.OpTodData:
		return nil, notImplErr
	case code.OpTodControl:
		return nil, notImplErr
	case code.OpRdm:
		return nil, notImplErr
	case code.OpRdmSub:
		return nil, notImplErr
	case code.OpMedia:
		return nil, notImplErr
	case code.OpMediaPatch:
		return nil, notImplErr
	case code.OpMediaControl:
		return nil, notImplErr
	case code.OpMediaContrlReply:
		return nil, notImplErr
	case code.OpTimeCode:
		p = &ArtTimeCodePacket{}
	case code.OpTimeSync:
		return nil, notImplErr
	case code.OpTrigger:
		p = &ArtTriggerPacket{}
	case code.OpDirectory:
		return nil, notImplErr
	case code.OpDirectoryReply:
		return nil, notImplErr
	case code.OpVideoSetup:
		return nil, notImplErr
	case code.OpVideoPalette:
		return nil, notImplErr
	case code.OpVideoData:
		return nil, notImplErr
	case code.OpMacMaster:
		return nil, notImplErr
	case code.OpMacSlave:
		return nil, notImplErr
	case code.OpFirmwareMaster:
		return nil, notImplErr
	case code.OpFirmwareReply:
		return nil, notImplErr
	case code.OpFileTnMaster:
		return nil, notImplErr
	case code.OpFileFnMaster:
		return nil, notImplErr
	case code.OpFileFnReply:
		return nil, notImplErr
	case code.OpIPProg:
		p = &ArtIPProgPacket{}
	case code.OpIPProgReply:
		p = &ArtIPProgReplyPacket{}
	default:
		return nil, notImplErr
	}

	err = p.UnmarshalBinary(b)
	return
}
