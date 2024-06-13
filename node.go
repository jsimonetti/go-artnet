package artnet

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

// Node is the information known about a node
type Node struct {
	// Config holds the configuration of this node
	Config NodeConfig

	log Logger

	// conn is the UDP connection this node will listen on
	conn          *net.UDPConn
	localAddr     net.UDPAddr
	broadcastAddr net.UDPAddr

	// sendCh is where payloads are passed to be sent by sendLoop
	sendCh chan netPayload
	// pollCh is how pollPackets are passed to the pollReplyLoop
	pollCh chan packet.ArtPollPacket

	packetHandlers map[code.OpCode]PacketHandler
}

// netPayload contains bytes read from the network and/or an error
type netPayload struct {
	address net.UDPAddr
	packet  packet.ArtNetPacket
}

// NewNode return a Node
func NewNode(name string, style code.StyleCode, ip net.IP, log Logger, opts ...NodeOption) *Node {
	n := &Node{
		Config: NodeConfig{
			Name: name,
			Type: style,
			IP:   ip,
		},
		log: log.With(Fields{"type": "Node"}),

		conn: nil,
		localAddr: net.UDPAddr{
			IP:   ip,
			Port: packet.ArtNetPort,
			Zone: "",
		},
		broadcastAddr: defaultBroadcastAddr,

		sendCh: make(chan netPayload, 10),
		pollCh: make(chan packet.ArtPollPacket, 10),

		packetHandlers: make(map[code.OpCode]PacketHandler),
	}

	if len(ip) < 1 {
		// TODO: generate an IP according to spec
		//ip = GenerateIP()
	}

	for _, opt := range opts {
		n.SetOption(opt)
	}

	handlePacketPoll := func(p packet.ArtNetPacket) {
		poll, ok := p.(*packet.ArtPollPacket)
		if !ok {
			n.log.With(Fields{"packet": p}).Debugf("unknown packet type")
			return
		}

		n.pollCh <- *poll
	}
	n.packetHandlers[code.OpPoll] = handlePacketPoll

	return n
}

// Start will start the controller.
// Closing ctx will end all routines and close connection
func (n *Node) Start(ctx context.Context) error {
	if err := n.Config.validate(); err != nil {
		return err
	}

	c, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", packet.ArtNetPort))
	if err != nil {
		n.log.With(Fields{"error": err}).Error("error net.ListenPacket")
		return err
	}
	n.conn = c.(*net.UDPConn)

	go n.pollReplyLoop(ctx)
	go n.recvLoop(ctx)
	go n.sendLoop(ctx)

	n.log.With(Fields{"ip": n.Config.IP.String(), "type": n.Config.Type.String()}).Debug("node started")

	return nil
}

// pollReplyLoop loops to reply to ArtPoll packets
// when a controller asks for continuous updates, we do that using a ticker
func (n *Node) pollReplyLoop(ctx context.Context) {
	var ticker time.Ticker

	// loop until shutdown
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			// if we should regularly send replies (can be requested by the controller)
			// we send it here

		case <-n.pollCh:
			// reply with pollReply
			// n.log.With(nil).Debug("sending ArtPollReply")

			n.sendCh <- netPayload{
				address: n.broadcastAddr,
				packet:  n.Config.buildArtPollReply(),
			}

			// TODO: if we are asked to send changes regularly, set the Ticker here

		}
	}
}

// sendLoop is used to send packets to the network
func (n *Node) sendLoop(ctx context.Context) {
	// loop until ctx ends
	for {
		select {
		case <-ctx.Done():
			err := n.conn.Close()
			if err != nil {
				n.log.With(Fields{"error": err}).Errorf("error closing conn")
			}
			return

		case payload := <-n.sendCh:

			// opCode := payload.packet.GetOpCode()
			// address := types.Address{}

			// if dmxPacket, ok := payload.packet.(*packet.ArtDMXPacket); ok {
			// 	address = dmxPacket.GetAddress()
			// }

			// create an ArtPoll packet to send out periodically
			b, err := payload.packet.MarshalBinary()
			if err != nil {
				n.log.With(Fields{"error": err, "dst": payload.address.String()}).Errorf("error marshalling packet")
				continue
			}
			_, err = n.conn.WriteToUDP(b, &payload.address)
			if err != nil {
				n.log.With(Fields{"error": err, "dst": payload.address.String()}).Errorf("error writing packet")
				continue
			}
			// n.log.With(Fields{
			// 	"dst":     payload.address.String(),
			// 	"bytes":   num,
			// 	"opCode":  opCode,
			// 	"address": address,
			// }).Debugf("packet sent")

		}
	}
}

// sendPacket will build and push a payload onto the sendCh
func (n *Node) send(a net.UDPAddr, p packet.ArtNetPacket) {

	// n.log.With(Fields{"dst": a.String(), "opCode": p.GetOpCode()}).Debugf("packet sending")

	n.sendCh <- netPayload{
		address: a,
		packet:  p,
	}
}

// recvLoop is used to receive packets from the network.
// due to the nature of broadcasting, we see our own sent
// packets to, but we ignore them
func (n *Node) recvLoop(ctx context.Context) {
	b := make([]byte, 4096)
	// loop until ctx ends
	for {
		select {
		case <-ctx.Done():
			// no need to close n.conn, sendLoop handles that
			return
		default:
		}

		num, from, err := n.conn.ReadFromUDP(b)

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

		// n.log.With(Fields{"src": from.String(), "bytes": num}).Debugf("received packet")

		p, err := packet.Unmarshal(b[:num])
		if err != nil {
			n.log.With(Fields{"data": fmt.Sprintf("%v", b)}).Warnf("failed to parse packet: %v", err)
			continue
		}
		go n.handlePacket(p)
	}
}

// handlePacket contains the logic for dealing with incoming packets
func (n *Node) handlePacket(p packet.ArtNetPacket) {
	handler, ok := n.packetHandlers[p.GetOpCode()]
	if !ok {
		n.log.With(Fields{"packet": p, "opCode": p.GetOpCode()}).Debugf("ignoring unhandled packet")
		return
	}

	go handler(p)
}

// PacketHandler gets called when a new packet has been received and needs to be processed
type PacketHandler func(p packet.ArtNetPacket)
