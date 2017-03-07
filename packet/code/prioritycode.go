//go:generate stringer -type=PriorityCode
package code

// PriorityCode defines Diagnostics Priority codes. These are used in ArtPoll and ArtDiagData
type PriorityCode uint8

const (
	// DpLow Low priority message.
	DpLow PriorityCode = 0x10

	// DpMed Medium priority message.
	DpMed PriorityCode = 0x40

	// DpHigh High priority message.
	DpHigh PriorityCode = 0x80

	// DpCritical Critical priority message.
	DpCritical PriorityCode = 0xe0

	// DpVolatile Volatile message. Messages of this type are displayed on a single line in the
	// DMX-Workshop diagnostics display. All other types are displayed in a list box.
	DpVolatile PriorityCode = 0xf0
)
