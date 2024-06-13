package artnet

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/types"
)

// ControlledNode holds the configuration of a node we control
type ControlledNode struct {
	nodeLock sync.Mutex

	BoundDevices map[types.BindIndex]NodeConfig

	LastSeen   time.Time
	UDPAddress net.UDPAddr

	Sequence uint8

	Outputs map[types.Address]bufferOutput
	Inputs  map[types.Address]InputPort
}

func newControlledNode(cfg NodeConfig) *ControlledNode {
	cn := &ControlledNode{
		BoundDevices: make(map[types.BindIndex]NodeConfig),

		LastSeen:   time.Now(),
		UDPAddress: net.UDPAddr{IP: cfg.IP, Port: packet.ArtNetPort},

		Sequence: 0,

		Outputs: make(map[types.Address]bufferOutput),
		Inputs:  make(map[types.Address]InputPort),
	}

	cn.update(cfg)
	return cn
}

func (cn *ControlledNode) update(cfg NodeConfig) {
	cn.nodeLock.Lock()
	defer cn.nodeLock.Unlock()

	cn.LastSeen = time.Now()

	if cfgPrev, ok := cn.BoundDevices[cfg.BindIndex]; ok {
		// previously known device

		// handle outputs
	NextOutput:
		for _, portPrev := range cfgPrev.OutputPorts {
			for _, port := range cfg.OutputPorts {
				if portPrev.Address == port.Address {
					// found matching port
					continue NextOutput
				}
			}
			delete(cn.Outputs, portPrev.Address)
		}

		for _, port := range cfg.OutputPorts {

			// If was not in DMXBuffers
			if _, ok := cn.Outputs[port.Address]; !ok {
				cn.Outputs[port.Address] = bufferOutput{
					&dmxBuffer{},
					port,
				}
			}
		}

		// handle inputs
		for _, port := range cfgPrev.InputPorts {
			delete(cn.Inputs, port.Address)
		}
		for _, port := range cfg.InputPorts {
			cn.Inputs[port.Address] = port
		}

	} else {
		// New device, just add ports

		for _, port := range cfg.OutputPorts {
			cn.Outputs[port.Address] = bufferOutput{
				&dmxBuffer{},
				port,
			}
		}
		for _, port := range cfg.InputPorts {
			cn.Inputs[port.Address] = port
		}
	}

	cn.BoundDevices[cfg.BindIndex] = cfg

}

func (cn *ControlledNode) RangeOutputs(f func(a types.Address)) {
	cn.nodeLock.Lock()

	addresses := make([]types.Address, len(cn.Outputs))

	i := 0
	for address := range cn.Outputs {
		addresses[i] = address
		i++
	}
	cn.nodeLock.Unlock()

	sortAddresses(addresses)

	for _, address := range addresses {
		f(address)
	}
}

func sortAddresses(addresses []types.Address) {
	sort.Slice(addresses, func(i, j int) bool {
		ii := addresses[i]
		jj := addresses[j]
		if ii.Net != jj.Net {
			return ii.Net < jj.Net
		}
		return ii.SubUni < jj.SubUni
	})

}

// getDMXUpdates will create all ArtDMXPackets which need to be sent
func (cn *ControlledNode) getDMXUpdates(activeInterval, passiveInterval time.Duration) []*packet.ArtDMXPacket {
	cn.nodeLock.Lock()
	defer cn.nodeLock.Unlock()

	packets := []*packet.ArtDMXPacket{}

	for address, output := range cn.Outputs {
		dmxBuffer := output.dmxBuffer
		dmx := dmxBuffer.checkUpdate(activeInterval, passiveInterval)
		if dmx == nil {
			continue
		}

		cn.Sequence++
		packet := packet.NewArtDMXPacket(address, *dmx, cn.Sequence)
		packets = append(packets, packet)
	}

	return packets
}

// setDMXBuffer will update the buffer on a universe address
func (cn *ControlledNode) SetDMXBuffer(address types.Address, dmx types.DMXData) error {
	cn.nodeLock.Lock()
	defer cn.nodeLock.Unlock()

	var output bufferOutput
	var ok bool

	if output, ok = cn.Outputs[address]; !ok {
		return fmt.Errorf("unknown address for controlled node")
	}

	output.set(dmx)

	return nil
}

type bufferOutput struct {
	*dmxBuffer
	OutputPort
}

type dmxBuffer struct {
	Port       OutputPort
	Data       types.DMXData
	LastUpdate time.Time
	Stale      bool
}

func (b *dmxBuffer) set(dmx types.DMXData) {
	b.Stale = true
	copy(b.Data[:], dmx[:])
}

func (b *dmxBuffer) checkUpdate(activeInterval, passiveInterval time.Duration) *types.DMXData {
	now := time.Now()
	durationSinceLastUpdate := now.Sub(b.LastUpdate)

	if b.Stale {
		if durationSinceLastUpdate < activeInterval {
			return nil
		}
	} else {
		if durationSinceLastUpdate < passiveInterval {
			return nil
		}
	}

	b.LastUpdate = now
	b.Stale = false

	var dmx types.DMXData
	copy(dmx[:], b.Data[:])
	return &dmx
}
