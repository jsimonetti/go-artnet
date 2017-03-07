package code

//go:generate stringer -type=OpCode

// OpCode defines the class of data following an UDP packet.
type OpCode uint16

// Valid returns wether the OpCode is valid
func Valid(o OpCode) bool {
	_, ok := _OpCode_map[o]
	return ok
}

const (
	// OpPoll This is an ArtPoll packet, no other data is contained in this UDP packet.
	OpPoll OpCode = 0x0020

	// OpPollReply This is an ArtPollReply Packet. It contains device status information.
	OpPollReply OpCode = 0x0021

	// OpDiagData Diagnostics and data logging packet.
	OpDiagData OpCode = 0x0023

	// OpCommand Used to send text based parameter commands.
	OpCommand OpCode = 0x0024

	// OpOutput This is an ArtDmx data packet. It contains zero start code DMX512 information for a single Universe.
	OpOutput OpCode = 0x0050

	// OpDMX This is an ArtDmx data packet. It contains zero start code DMX512 information for a single Universe.
	OpDMX OpCode = OpOutput

	// OpNzs This is an ArtNzs data packet. It contains non-zero start code (except RDM) DMX512 information for a single Universe.
	OpNzs OpCode = 0x0051

	// OpSync This is an ArtSync data packet. It is used to force synchronous transfer of ArtDmx packets to a node’s output.
	OpSync OpCode = 0x0052

	// OpAddress This is an ArtAddress packet. It contains remote programming information for a Node.
	OpAddress OpCode = 0x0060

	// OpInput This is an ArtInput packet. It contains enable–disable data for DMX inputs.
	OpInput OpCode = 0x0070

	// OpTodRequest This is an ArtTodRequest packet. It is used to request a Table of Devices (ToD) for RDM discovery.
	OpTodRequest OpCode = 0x0080

	// OpTodData This is an ArtTodData packet. It is used to send a Table of Devices (ToD) for RDM discovery.
	OpTodData OpCode = 0x0081

	// OpTodControl This is an ArtTodControl packet. It is used to send RDM discovery control messages.
	OpTodControl OpCode = 0x0082

	// OpRdm This is an ArtRdm packet. It is used to send all non discovery RDM messages.
	OpRdm OpCode = 0x0083

	// OpRdmSub This is an ArtRdmSub packet. It is used to send compressed, RDM Sub-Device data.
	OpRdmSub OpCode = 0x0084

	// OpMedia This is an ArtMedia packet. It is Unicast by a Media Server and acted upon by a Controller.
	OpMedia OpCode = 0x0090

	// OpMediaPatch This is an ArtMediaPatch packet. It is Unicast by a Controller and acted upon by a Media Server.
	OpMediaPatch OpCode = 0x0091

	// OpMediaControl This is an ArtMediaControl packet. It is Unicast by a Controller and acted upon by a Media Server.
	OpMediaControl OpCode = 0x0092

	// OpMediaContrlReply This is an ArtMediaControlReply packet. It is Unicast by a Media Server and acted upon by a Controller.
	OpMediaContrlReply OpCode = 0x0093

	// OpTimeCode This is an ArtTimeCode packet. It is used to transport time code over the network.
	OpTimeCode OpCode = 0x0097

	// OpTimeSync Used to synchronise real time date and clock
	OpTimeSync OpCode = 0x0098

	// OpTrigger Used to send trigger macros
	OpTrigger OpCode = 0x0099

	// OpDirectory Requests a node's file list
	OpDirectory OpCode = 0x009a

	// OpDirectoryReply Replies to OpDirectory with file list
	OpDirectoryReply OpCode = 0x009b

	// OpVideoSetup This is an ArtVideoSetup packet. It contains video screen setup information for nodes that implement the extended video features.
	OpVideoSetup OpCode = 0x10a0

	// OpVideoPalette This is an ArtVideoPalette packet. It contains colour palette setup information for nodes that implement the extended video features.
	OpVideoPalette OpCode = 0x20a0

	// OpVideoData This is an ArtVideoData packet. It contains display data for nodes that implement the extended video features.
	OpVideoData OpCode = 0x40a0

	// OpMacMaster This packet is deprecated.
	OpMacMaster OpCode = 0x00f0

	// OpMacSlave This packet is deprecated.
	OpMacSlave OpCode = 0x00f1

	// OpFirmwareMaster This is an ArtFirmwareMaster packet. It is used to upload new firmware or firmware extensions to theNode.
	OpFirmwareMaster OpCode = 0x00f2

	// OpFirmwareReply This is an ArtFirmwareReply packet. It is returned by the node to acknowledge receipt of an ArtFirmwareMaster packet or ArtFileTnMaster packet.
	OpFirmwareReply OpCode = 0x00f3

	// OpFileTnMaster Uploads user file to node.
	OpFileTnMaster OpCode = 0x00f4

	// OpFileFnMaster Downloads user file from node.
	OpFileFnMaster OpCode = 0x00f5

	// OpFileFnReply Server to Node acknowledge for download packets.
	OpFileFnReply OpCode = 0x00f6

	// OpIPProg This is an ArtIpProg packet. It is used to re-programme the IP address andMask of the Node.
	OpIPProg OpCode = 0x00f8

	// OpIPProgReply This is an ArtIpProgReply packet. It is returned by the node to acknowledge receipt of an ArtIpProg packet.
	OpIPProgReply OpCode = 0x00f9
)
