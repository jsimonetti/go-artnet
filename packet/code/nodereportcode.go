package code

//go:generate stringer -type=NodeReportCode

// NodeReportCode defines generic error, advisory and status messages for both Nodes and Controllers
type NodeReportCode uint8

const (
	// RcDebug Booted in debug mode (Only used in development)
	RcDebug NodeReportCode = 0x00

	// RcPowerOk Power On Tests successful
	RcPowerOk NodeReportCode = 0x01

	// RcPowerFail Hardware tests failed at Power On
	RcPowerFail NodeReportCode = 0x02

	// RcSocketWr1 Last UDP from Node failed due to truncated length, Most likely caused by a collision.
	RcSocketWr1 NodeReportCode = 0x03

	// RcParseFail Unable to identify last UDP transmission. Check OpCode and packet length.
	RcParseFail NodeReportCode = 0x04

	// RcUDPFail Unable to open Udp Socket in last transmission attempt
	RcUDPFail NodeReportCode = 0x05

	// RcShNameOk Confirms that Short Name programming via ArtAddress, was successful.
	RcShNameOk NodeReportCode = 0x06

	// RcLoNameOk Confirms that Long Name programming via ArtAddress, was successful.
	RcLoNameOk NodeReportCode = 0x07

	// RcDmxError DMX512 receive errors detected.
	RcDmxError NodeReportCode = 0x08

	// RcDmxUDPFull Ran out of internal DMX transmit buffers.
	RcDmxUDPFull NodeReportCode = 0x09

	// RcDmxRxFull Ran out of internal DMX Rx buffers.
	RcDmxRxFull NodeReportCode = 0x0a

	// RcSwitchErr Rx Universe switches conflict.
	RcSwitchErr NodeReportCode = 0x0b

	// RcConfigErr Product configuration does not match firmware.
	RcConfigErr NodeReportCode = 0x0c

	// RcDmxShort DMX output short detected. See GoodOutput field.
	RcDmxShort NodeReportCode = 0x0d

	// RcFirmwareFail Last attempt to upload new firmware failed.
	RcFirmwareFail NodeReportCode = 0x0e

	// RcUserFail User changed switch settings when address locked by remote programming. User changes ignored.
	RcUserFail NodeReportCode = 0x0f

	// RcFactoryRes Factory reset has occurred
	RcFactoryRes NodeReportCode = 0x10
)
