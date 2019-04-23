package artnet

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

// Node is the information known about a node
type Node struct {
	// Config holds the configuration of this node
	Config NodeConfig

	// conn is the UDP connection this node will listen on
	conn      *net.UDPConn
	localAddr net.UDPAddr
	sendCh    chan netPayload
	recvCh    chan netPayload

	// shutdownCh will be closed on shutdown of the node
	shutdownCh   chan struct{}
	shutdown     bool
	shutdownErr  error
	shutdownLock sync.Mutex

	// pollCh will receive ArtPoll packets
	pollCh chan packet.ArtPollPacket
	// pollCh will send ArtPollReply packets
	pollReplyCh chan packet.ArtPollReplyPacket

	log Logger
}

// netPayload contains bytes read from the network and/or an error
type netPayload struct {
	address net.UDPAddr
	err     error
	data    []byte
}

// NewNode return a Node
func NewNode(name string, style code.StyleCode, ip net.IP, log Logger) *Node {
	n := &Node{
		Config: NodeConfig{
			Name: name,
			Type: style,
		},
		conn:     nil,
		shutdown: true,
		log:      log.With(Fields{"type": "Node"}),
	}

	if len(ip) < 1 {
		// TODO: generate an IP according to spec
		//ip = GenerateIP()
	}
	n.Config.IP = ip
	n.localAddr = net.UDPAddr{
		IP:   ip,
		Port: packet.ArtNetPort,
		Zone: "",
	}

	return n
}

// Stop will stop all running routines and close the network connection
func (n *Node) Stop() {
	n.shutdownLock.Lock()
	n.shutdown = true
	n.shutdownLock.Unlock()
	close(n.shutdownCh)
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			n.log.Printf("failed to close read socket: %v")
		}
	}
}

func (n *Node) isShutdown() bool {
	n.shutdownLock.Lock()
	defer n.shutdownLock.Unlock()
	return n.shutdown
}

// Start will start the controller
func (n *Node) Start() error {
	n.log.With(Fields{"ip": n.Config.IP.String(), "type": n.Config.Type.String()}).Debug("node started")

	n.sendCh = make(chan netPayload, 10)
	n.recvCh = make(chan netPayload, 10)
	n.pollCh = make(chan packet.ArtPollPacket, 10)
	n.pollReplyCh = make(chan packet.ArtPollReplyPacket, 10)
	n.shutdownCh = make(chan struct{})
	n.shutdown = false

	c, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", packet.ArtNetPort))
	if err != nil {
		n.shutdownErr = fmt.Errorf("error net.ListenPacket: %s", err)
		n.log.With(Fields{"error": err}).Error("error net.ListenPacket")
		return err
	}
	n.conn = c.(*net.UDPConn)

	go n.pollReplyLoop()
	go n.recvLoop()
	go n.sendLoop()

	return nil
}

// pollReplyLoop loops to reply to ArtPoll packets
// when a controller asks for continuous updates, we do that using a ticker
func (n *Node) pollReplyLoop() {
	var timer time.Ticker

	// create an ArtPollReply packet to send out in response to an ArtPoll packet
	p := ArtPollReplyFromConfig(n.Config)
	me, err := p.MarshalBinary()
	if err != nil {
		n.log.With(Fields{"err": err}).Error("error creating ArtPollReply packet for self")
		return
	}

	// loop until shutdown
	for {
		select {
		case <-timer.C:
			// if we should regularly send replies (can be requested by the controller)
			// we send it here

		case <-n.pollCh:
			// reply with pollReply
			n.log.With(nil).Debug("poll received, now send a reply")

			n.sendCh <- netPayload{
				address: broadcastAddr,
				data:    me,
			}

			// TODO: if we are asked to send changes regularly, set the Ticker here

		case <-n.shutdownCh:
			return
		}
	}
}

// sendLoop is used to send packets to the network
func (n *Node) sendLoop() {
	// loop until shutdown
	for {
		select {
		case <-n.shutdownCh:
			return

		case payload := <-n.sendCh:
			if n.isShutdown() {
				return
			}

			num, err := n.conn.WriteToUDP(payload.data, &payload.address)
			if err != nil {
				n.log.With(Fields{"error": err}).Debugf("error writing packet")
				continue
			}
			n.log.With(Fields{"dst": payload.address.String(), "bytes": num}).Debugf("packet sent")

		}
	}
}

// recvLoop is used to receive packets from the network
// it starts a goroutine for dumping the msgs onto a channel,
// the payload from that channel is then fed into a handler
// due to the nature of broadcasting, we see our own sent
// packets to, but we ignore them
func (n *Node) recvLoop() {
	// start a routine that will read data from n.conn
	// and (if not shutdown), send to the recvCh
	go func() {
		b := make([]byte, 4096)
		for {
			num, from, err := n.conn.ReadFromUDP(b)
			if n.isShutdown() {
				return
			}

			if n.localAddr.IP.Equal(from.IP) {
				// this was sent by me, so we ignore it
				//n.log.With(Fields{"src": from.String(), "bytes": num}).Debugf("ignoring received packet from self")
				continue
			}

			if err != nil {
				if err == io.EOF {
					return
				}

				n.log.With(Fields{"src": from.String(), "bytes": num}).Errorf("failed to read from socket: %v", err)
				continue
			}

			n.log.With(Fields{"src": from.String(), "bytes": num}).Debugf("received packet")
			payload := netPayload{
				address: *from,
				err:     err,
				data:    make([]byte, num),
			}
			copy(payload.data, b)
			n.recvCh <- payload
		}
	}()

	// loop until shutdown
	for {
		select {
		case payload := <-n.recvCh:
			p, err := packet.Unmarshal(payload.data)
			if err != nil {
				n.log.With(Fields{
					"src":  payload.address.IP.String(),
					"data": fmt.Sprintf("%v", payload.data),
				}).Warnf("failed to parse packet: %v", err)
				continue
			}
			go n.handlePacket(p)

		case <-n.shutdownCh:
			return
		}
	}
}

// handlePacket contains the logic for dealing with incoming packets
func (n *Node) handlePacket(p packet.ArtNetPacket) {
	switch p := p.(type) {
	case *packet.ArtPollReplyPacket:
		// only handle these packets if we are a controller
		if n.Config.Type == code.StController {
			n.pollReplyCh <- *p
		}

	case *packet.ArtPollPacket:
		n.pollCh <- *p

	default:
		n.log.With(Fields{"packet": p}).Debugf("unknown packet type")
	}

}
