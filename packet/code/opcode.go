package code

//go:generate stringer -type=OpCode

import "fmt"

// OpCode defines the class of data following an UDP packet.
type OpCode uint16

// Unmarshal unmarshals the contents of a byte slice into an OpCode.
func (o OpCode) Unmarshal(b []byte) OpCode {
	if len(b) != 2 {
		fmt.Printf("len != 2: %d", len(b))
		return OpCode(0)
	}
	//Opcode is always transmitted low byte first
	c := OpCode(uint16(b[0]) | uint16(b[1])<<8)
	return c
}

// Marshal marshals an OpCode into a byte slice.
func (o *OpCode) Marshal() []byte {
	//Opcode is always transmitted low byte first
	b := make([]byte, 2)
	b[0] = uint8(*o)
	b[1] = uint8(*o >> 8)
	return b
}

// Valid returns wether the OpCode is valid
func (o *OpCode) Valid() bool {
	_, ok := _OpCode_map[*o]
	return ok
}

const (
	// OpPoll This is an ArtPoll packet, no other data is contained in this UDP packet.
	OpPoll OpCode = 0x2000

	// OpPollReply This is an ArtPollReply Packet. It contains device status information.
	OpPollReply OpCode = 0x2100

	// OpDiagData Diagnostics and data logging packet.
	OpDiagData OpCode = 0x2300

	// OpCommand Used to send text based parameter commands.
	OpCommand OpCode = 0x2400

	// OpOutput This is an ArtDmx data packet. It contains zero start code DMX512 information for a single Universe.
	OpOutput OpCode = 0x5000

	// OpDMX This is an ArtDmx data packet. It contains zero start code DMX512 information for a single Universe.
	OpDMX OpCode = OpOutput

	// OpNzs This is an ArtNzs data packet. It contains non-zero start code (except RDM) DMX512 information for a single Universe.
	OpNzs OpCode = 0x5100

	// OpSync This is an ArtSync data packet. It is used to force synchronous transfer of ArtDmx packets to a node’s output.
	OpSync OpCode = 0x5200

	// OpAddress This is an ArtAddress packet. It contains remote programming information for a Node.
	OpAddress OpCode = 0x6000

	// OpInput This is an ArtInput packet. It contains enable–disable data for DMX inputs.
	OpInput OpCode = 0x7000

	// OpTodRequest This is an ArtTodRequest packet. It is used to request a Table of Devices (ToD) for RDM discovery.
	OpTodRequest OpCode = 0x8000

	// OpTodData This is an ArtTodData packet. It is used to send a Table of Devices (ToD) for RDM discovery.
	OpTodData OpCode = 0x8100

	// OpTodControl This is an ArtTodControl packet. It is used to send RDM discovery control messages.
	OpTodControl OpCode = 0x8200

	// OpRdm This is an ArtRdm packet. It is used to send all non discovery RDM messages.
	OpRdm OpCode = 0x8300

	// OpRdmSub This is an ArtRdmSub packet. It is used to send compressed, RDM Sub-Device data.
	OpRdmSub OpCode = 0x8400

	// OpMedia This is an ArtMedia packet. It is Unicast by a Media Server and acted upon by a Controller.
	OpMedia OpCode = 0x9000

	// OpMediaPatch This is an ArtMediaPatch packet. It is Unicast by a Controller and acted upon by a Media Server.
	OpMediaPatch OpCode = 0x9100

	// OpMediaControl This is an ArtMediaControl packet. It is Unicast by a Controller and acted upon by a Media Server.
	OpMediaControl OpCode = 0x9200

	// OpMediaContrlReply This is an ArtMediaControlReply packet. It is Unicast by a Media Server and acted upon by a Controller.
	OpMediaContrlReply OpCode = 0x9300

	// OpTimeCode This is an ArtTimeCode packet. It is used to transport time code over the network.
	OpTimeCode OpCode = 0x9700

	// OpTimeSync Used to synchronise real time date and clock
	OpTimeSync OpCode = 0x9800

	// OpTrigger Used to send trigger macros
	OpTrigger OpCode = 0x9900

	// OpDirectory Requests a node's file list
	OpDirectory OpCode = 0x9a00

	// OpDirectoryReply Replies to OpDirectory with file list
	OpDirectoryReply OpCode = 0x9b00

	// OpVideoSetup This is an ArtVideoSetup packet. It contains video screen setup information for nodes that implement the extended video features.
	OpVideoSetup OpCode = 0xa010

	// OpVideoPalette This is an ArtVideoPalette packet. It contains colour palette setup information for nodes that implement the extended video features.
	OpVideoPalette OpCode = 0xa020

	// OpVideoData This is an ArtVideoData packet. It contains display data for nodes that implement the extended video features.
	OpVideoData OpCode = 0xa040

	// OpMacMaster This packet is deprecated.
	OpMacMaster OpCode = 0xf000

	// OpMacSlave This packet is deprecated.
	OpMacSlave OpCode = 0xf100

	// OpFirmwareMaster This is an ArtFirmwareMaster packet. It is used to upload new firmware or firmware extensions to theNode.
	OpFirmwareMaster OpCode = 0xf200

	// OpFirmwareReply This is an ArtFirmwareReply packet. It is returned by the node to acknowledge receipt of an ArtFirmwareMaster packet or ArtFileTnMaster packet.
	OpFirmwareReply OpCode = 0xf300

	// OpFileTnMaster Uploads user file to node.
	OpFileTnMaster OpCode = 0xf400

	// OpFileFnMaster Downloads user file from node.
	OpFileFnMaster OpCode = 0xf500

	// OpFileFnReply Server to Node acknowledge for download packets.
	OpFileFnReply OpCode = 0xf600

	// OpIPProg This is an ArtIpProg packet. It is used to re-programme the IP address andMask of the Node.
	OpIPProg OpCode = 0xf800

	// OpIPProgReply This is an ArtIpProgReply packet. It is returned by the node to acknowledge receipt of an ArtIpProg packet.
	OpIPProgReply OpCode = 0xf900
)
