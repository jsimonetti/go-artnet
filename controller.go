package artnet

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/jsimonetti/go-artnet/artnettypes"
	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

var defaultBroadcastAddr = net.UDPAddr{
	IP:   []byte{0x02, 0xff, 0xff, 0xff},
	Port: int(packet.ArtNetPort),
}

// Controller holds the information for a controller
type Controller struct {
	// cNode is the Node for the cNode
	cNode *Node

	log           Logger
	broadcastAddr net.UDPAddr

	// nodes is map of IPs to ControlledNodes that are seen by this controller
	nodes    map[string]*ControlledNode
	nodeLock sync.Mutex

	updateInterval       time.Duration
	expectActiveInterval time.Duration
	pollInterval         time.Duration
}

// NewController return a Controller
func NewController(name string, ip net.IP, log Logger, opts ...Option) *Controller {
	c := &Controller{
		log:           log,
		broadcastAddr: defaultBroadcastAddr,

		nodes: map[string]*ControlledNode{},

		updateInterval:       time.Millisecond * 30,
		expectActiveInterval: time.Second * 10, // nodes are stale after 5 missed ArtPoll's
		pollInterval:         time.Second * 2,
	}

	handlePacketPollReply := func(p packet.ArtNetPacket) {
		pollReply, ok := p.(*packet.ArtPollReplyPacket)
		if !ok {
			c.log.With(Fields{"packet": p}).Debugf("unknown packet type")
			return
		}

		cfg := newNodeConfigFrom(pollReply)

		switch cfg.Type {
		case code.StController:
			if len(cfg.OutputPorts) == 0 {
				// we don't care for controllers which do not have output ports for now // @todo
				// otherwise we simply treat controllers like nodes unless controller to controller
				// communication is implemented according to Art-Net specification
				return
			}
			if cfg.IP.String() == c.cNode.Config.IP.String() {
				// we don't care to keep track of ourself
				return
			}
		case code.StNode:
			break
		default:
			// TODO handle other types of devices
			return
		}

		c.updateNode(cfg)
	}

	c.cNode = NewNode(
		name,
		code.StController,
		ip,
		log,
		NodeOptionPacketHandler(code.OpPollReply, handlePacketPollReply),
	)

	for _, opt := range opts {
		c.SetOption(opt)
	}

	return c
}

// Start will start this controller
// Closing ctx will end all routines and close connection
func (c *Controller) Start(ctx context.Context) error {
	c.log = c.log.With(Fields{"type": "Controller"})

	if err := c.cNode.Start(ctx); err != nil {
		return fmt.Errorf("failed to start controller node: %v", err)
	}

	go c.pollLoop(ctx)
	go c.dmxUpdateLoop(ctx)
	return nil
}

// pollLoop will routinely poll for new nodes
func (c *Controller) pollLoop(ctx context.Context) {
	pollTicker := time.NewTicker(c.pollInterval)

	artPoll := packet.NewArtPollPacket()
	artPoll.TalkToMe = new(code.TalkToMe).WithReplyOnChange(true)
	artPoll.Priority = code.DpAll

	c.cNode.send(c.broadcastAddr, artPoll)

	// loop until shutdown
	for {
		select {
		case <-ctx.Done():
			return

		case <-pollTicker.C:
			c.cNode.send(c.broadcastAddr, artPoll)
		}
	}
}

// SendDMX will set the DMXBuffer for a destination address
func (c *Controller) SendDMX(ip net.IP, address artnettypes.Address, dmx [512]byte) error {
	// c.log.With(Fields{"address": address.String()}).Debug("received update channels")

	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	cn, ok := c.nodes[ip.String()]
	if !ok {
		err := errors.New("could not find node for ip")
		c.log.With(Fields{"ip": ip.String()}).Error(err)
		return err
	}

	err := cn.SetDMXBuffer(address, dmx)
	if err != nil {
		c.log.With(Fields{"err": err, "address": address.String()}).Error(err)
		return err
	}

	return nil
}

// dmxUpdateLoop will periodically update nodes until ctx ends
func (c *Controller) dmxUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(c.updateInterval)

	// loop until shutdown
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			c.dmxUpdate()
		}
	}
}

// dmxUpdate will updates nodes
func (c *Controller) dmxUpdate() {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	// send DMX buffer update
	for _, node := range c.nodes {
		packets := node.getDMXUpdates()
		for _, packet := range packets {
			c.cNode.send(node.udpAddress, packet)
		}
	}
	c.cNode.send(c.broadcastAddr, packet.NewArtSyncPacket())
}

// updateNode will either update an existing node or create a new one if the cfg has a unique IP
func (c *Controller) updateNode(cfg NodeConfig) {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	node, ok := c.nodes[cfg.IP.String()]
	if ok {
		// c.log.With(Fields{"node": cfg.Name, "ip": cfg.IP.String()}).Debug("updated node")
		// c.LogNode(cfg)
		node.update(cfg)
		return
	}

	c.nodes[cfg.IP.String()] = newControlledNode(cfg)
}

// expireNodesLoop will remove stale Nodes from the list of known nodes
// it will loop through the list of nodes and remove nodes older then X seconds
func (c *Controller) expireNodesLoop(ctx context.Context) {
	ticker := time.NewTicker(c.expectActiveInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.expireNodes()
		}
	}
}

func (c *Controller) expireNodes() {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for ip, node := range c.nodes {
		if time.Now().Sub(node.lastSeen) > c.expectActiveInterval {
			// node is stale
			c.log.With(Fields{"ip": ip}).Debug("removing stale node")
			delete(c.nodes, ip)
		}
	}
}

func (c *Controller) GetNode(ip net.IP) (*ControlledNode, error) {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	cn, ok := c.nodes[ip.String()]
	if !ok {
		err := errors.New("could not find node for ip")
		c.log.With(Fields{"ip": ip.String()}).Error(err)
		return nil, err
	}
	return cn, nil
}

func (c *Controller) LogNodes() {
	c.log.With(nil).Info("Logging known nodes...\n")

	c.RangeNodes(func(cn *ControlledNode) {
		for _, cfg := range cn.BoundDevices {
			c.LogNode(cfg)
		}
	})

	c.log.With(nil).Info("Logged known nodes.\n")
}

func (c *Controller) LogNode(cfg NodeConfig) {
	outs := make([]artnettypes.Address, len(cfg.OutputPorts))
	for i, port := range cfg.OutputPorts {
		outs[i] = port.Address
	}
	sortAddresses(outs)

	ins := make([]artnettypes.Address, len(cfg.InputPorts))
	for i, port := range cfg.InputPorts {
		ins[i] = port.Address
	}
	sortAddresses(ins)

	c.log.With(Fields{
		// "name": cfg.Name,
		// "description":  cfg.Description,
		// "manufacturer": cfg.Manufacturer,
		"ip": cfg.IP.String(),
		// "port": cfg.Port,
		// "hardwareAddr": cfg.Ethernet,
		// "bindIP":       cfg.BindIP,
		"bindIndex": cfg.BindIndex,
		// "report":    cfg.Report,
		// "status1":      cfg.Status1,
		// "status2":      cfg.Status2,
		// "status3":      cfg.Status3,
		"baseAddress": cfg.BaseAddress,
		"outPorts":    outs,
		"inPorts":     ins,
	}).Info("Logging Node Data")
}

func (c *Controller) RangeNodes(f func(cn *ControlledNode)) {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	for _, cn := range c.nodes {
		f(cn)
	}
}

func (c *Controller) RangeAll(f func(ip string, a artnettypes.Address)) {
	c.rangeIPs(func(ip string) {
		c.rangeOutputsOf(ip, func(a artnettypes.Address) {
			f(ip, a)
		})
	})
}

func (c *Controller) rangeIPs(f func(ip string)) {
	c.nodeLock.Lock()

	ips := make([]string, len(c.nodes))

	i := 0
	for ip := range c.nodes {
		ips[i] = ip
		i++
	}
	c.nodeLock.Unlock()

	sort.Strings(ips)

	for _, ip := range ips {
		f(ip)
	}
}

func (c *Controller) rangeOutputsOf(ip string, f func(a artnettypes.Address)) {
	c.nodeLock.Lock()
	defer c.nodeLock.Unlock()

	if cn, ok := c.nodes[ip]; ok {
		go cn.RangeOutputs(f)
	}
}
