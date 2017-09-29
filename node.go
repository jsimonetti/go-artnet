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
	conn   net.PacketConn
	bconn  net.Conn
	sendCh chan netPayload
	recvCh chan netPayload

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
	if len(ip) > 0 {
		n.Config.IP = ip
	}
	//n.Config.IP = GenerateIP()
	return n
}

// Stop will stop all running routines and close the network connection
func (n *Node) Stop() {
	n.shutdownLock.Lock()
	n.shutdown = true
	n.shutdownLock.Unlock()
	close(n.shutdownCh)
	if n.conn != nil {
		n.conn.Close()
	}
	if n.bconn != nil {
		n.bconn.Close()
	}
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

	var err error

	n.bconn, err = net.Dial("udp4", "2.255.255.255:6454")
	if err != nil {
		n.shutdownErr = fmt.Errorf("error net.ListenUDP: %s", err)
		n.log.With(Fields{"error": err}).Error("error net.ListenUDP")
		return err
	}

	n.conn, err = net.ListenPacket("udp4", "0.0.0.0:6454")
	if err != nil {
		n.shutdownErr = fmt.Errorf("error net.ListenUDP: %s", err)
		n.log.With(Fields{"error": err}).Error("error net.ListenUDP")
		return err
	}

	go n.recvLoop()
	go n.sendLoop()

	return nil
}

// pollReplyLoop loops to reply to ArtPoll packets
// when a controller asks for continuous updates, we do that using a ticker
func (n *Node) pollReplyLoop() {
	var timer time.Ticker

	// loop untill shutdown
	for {
		select {
		case <-timer.C:
			// if we should regularly send replies (can be requested by the controller)
			// we send it here

		case poll := <-n.pollCh:
			// reply with pollReply
			n.log.With(Fields{"poll": poll}).Debugf("poll received, now send a reply")

			// if we are asked to send changes regularyl, set the Ticker here

		case <-n.shutdownCh:
			return
		}
	}
}

// sendLoop is used to send packets to the network
func (n *Node) sendLoop() {
	// loop untill shutdown
	for {
		select {
		case payload := <-n.sendCh:
			n.shutdownLock.Lock()
			if !n.shutdown {
				var num int
				var err error
				if payload.address.IP.Equal(broadcastAddr.IP) {
					num, err = n.bconn.Write(payload.data)
				} else {
					num, err = n.conn.WriteTo(payload.data, &payload.address)
				}
				if err != nil {
					n.log.With(Fields{"error": err}).Debugf("error writing packet")
					continue
				}
				n.log.With(Fields{"dst": payload.address.String(), "bytes": num}).Debugf("packet sent")
			}
			n.shutdownLock.Unlock()
		case <-n.shutdownCh:
			return
		}
	}
}

// AddrToUDPAddr will turn a net.Addr into a net.UDPAddr
func AddrToUDPAddr(addr net.Addr) net.UDPAddr {
	udp := addr.(*net.UDPAddr)
	return *udp
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
			num, src, err := n.conn.ReadFrom(b)
			n.shutdownLock.Lock()
			if !n.shutdown {
				n.shutdownLock.Unlock()
				from := AddrToUDPAddr(src)
				if n.Config.IP.Equal(from.IP) {
					// this was sent by me, so we ignore it
					//n.log.With(Fields{"src": from.String(), "bytes": num}).Debugf("ignoring received packet from self")
					continue
				}

				n.log.With(Fields{"src": from.String(), "bytes": num}).Debugf("received packet")
				if err != nil && err != io.EOF {
					n.recvCh <- netPayload{
						address: from,
						data:    b[:num],
						err:     err,
					}
					continue
				}
				n.recvCh <- netPayload{
					address: from,
					data:    b[:num],
					err:     err,
				}
				continue
			}
			n.shutdownLock.Unlock()
			return
		}
	}()

	// loop untill shutdown
	for {
		select {
		case payload := <-n.recvCh:
			//if payload.err == nil {
			p, err := packet.Unmarshal(payload.data)
			if err == nil {
				// if this is a valid packet we handle it
				go n.handlePacket(p)
			}
			//}

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

	default:
		n.log.With(Fields{"packet": p}).Debugf("unknown packet type")
	}

}
