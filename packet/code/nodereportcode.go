//go:generate stringer -type=NodeReportCode
package code

// NodeReportCode defines generic error, advisory and status messages for both Nodes and Controllers
type NodeReportCode uint8

const (
	// RcDebug Booted in debug mode (Only used in development)
	RcDebug NodeReportCode = 0x0000

	// RcPowerOk Power On Tests successful
	RcPowerOk NodeReportCode = 0x0001

	// RcPowerFail Hardware tests failed at Power On
	RcPowerFail NodeReportCode = 0x0002

	// RcSocketWr1 Last UDP from Node failed due to truncated length, Most likely caused by a collision.
	RcSocketWr1 NodeReportCode = 0x0003

	// RcParseFail Unable to identify last UDP transmission. Check OpCode and packet length.
	RcParseFail NodeReportCode = 0x0004

	// RcUDPFail Unable to open Udp Socket in last transmission attempt
	RcUDPFail NodeReportCode = 0x0005

	// RcShNameOk Confirms that Short Name programming via ArtAddress, was successful.
	RcShNameOk NodeReportCode = 0x0006

	// RcLoNameOk Confirms that Long Name programming via ArtAddress, was successful.
	RcLoNameOk NodeReportCode = 0x0007

	// RcDmxError DMX512 receive errors detected.
	RcDmxError NodeReportCode = 0x0008

	// RcDmxUDPFull Ran out of internal DMX transmit buffers.
	RcDmxUDPFull NodeReportCode = 0x0009

	// RcDmxRxFull Ran out of internal DMX Rx buffers.
	RcDmxRxFull NodeReportCode = 0x000a

	// RcSwitchErr Rx Universe switches conflict.
	RcSwitchErr NodeReportCode = 0x000b

	// RcConfigErr Product configuration does not match firmware.
	RcConfigErr NodeReportCode = 0x000c

	// RcDmxShort DMX output short detected. See GoodOutput field.
	RcDmxShort NodeReportCode = 0x000d

	// RcFirmwareFail Last attempt to upload new firmware failed.
	RcFirmwareFail NodeReportCode = 0x000e

	// RcUserFail User changed switch settings when address locked by remote programming. User changes ignored.
	RcUserFail NodeReportCode = 0x000f

	// RcFactoryRes Factory reset has occurred
	RcFactoryRes NodeReportCode = 0x0010
)
