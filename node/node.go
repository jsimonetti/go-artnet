package node

import (
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
	Config Config

	// conn is the UDP connection this node will listen on
	conn   *net.UDPConn
	sendCh chan *netPayload
	recvCh chan *netPayload

	// shutdownCh will be closed on shutdown of the node
	shutdownCh chan struct{}
	shutdown   bool

	// pollCh will receive ArtPoll packets
	pollCh chan *packet.ArtPollPacket

	// Controller is a config of a controller should this node by under it's controller
	Controller Config
}

type netPayload struct {
	err  error
	data []byte
}

// New return a Node
func New(name string, style code.StyleCode, ip net.IP) Node {
	n := Node{
		Config: Config{
			Name: name,
			Type: style,
		},
		conn:     nil,
		shutdown: true,
	}
	if len(ip) > 0 {
		n.Config.IP = ip
	}
	//n.Config.IP = GenerateIP()
	return n
}

// Stop will stop all running routines and close the network connection
func (n *Node) Stop() {
	n.shutdown = true
	close(n.shutdownCh)
	if n.conn != nil {
		n.conn.Close()
		n.conn = nil
	}
	close(n.sendCh)
	close(n.recvCh)
	close(n.pollCh)
}

// Start will start the controller
func (n *Node) Start() (err error) {
	n.sendCh = make(chan *netPayload)
	n.recvCh = make(chan *netPayload)
	n.pollCh = make(chan *packet.ArtPollPacket)
	n.shutdownCh = make(chan struct{})
	n.shutdown = false

	src := fmt.Sprintf("%s:%d", n.Config.IP, packet.ArtNetPort)
	localAddr, _ := net.ResolveUDPAddr("udp", src)

	n.conn, err = net.ListenUDP("udp", localAddr)
	if err != nil {
		return fmt.Errorf("error net.ListenUDP: %s", err)
	}

	go n.recvLoop()
	go n.sendLoop()

	select {
	case <-n.shutdownCh:
		return nil
	}
}

func (n *Node) pollReplyLoop() {
	var timer time.Ticker

	for {
		select {
		case <-timer.C:
			// if we should regularly send replies (can be requested by the controller)
			// we send it here

		case poll := <-n.pollCh:
			// reply with pollReply
			fmt.Printf("poll received: %v, now send a reply", poll)

			// if we are asked to send changes regularyl, set the Ticker here

		case <-n.shutdownCh:
			return
		}
	}
}

func (n *Node) sendLoop() {
	for {
		select {
		case p := <-n.sendCh:
			fmt.Printf("send %v", p)
		case <-n.shutdownCh:
			return
		}
	}
}

func (n *Node) recvLoop() {
	// start a routine that will read data from n.conn
	// and (if not shutdown), send to the recvCh
	go func() {
		b := make([]byte, 4096)
		for {
			num, from, err := n.conn.ReadFromUDP(b)
			if !n.shutdown {
				if from.IP.Equal(n.Config.IP) {
					// this was sent from me, so we ignore it
					continue
				}
				if err != nil && err != io.EOF {
					n.recvCh <- &netPayload{
						data: b[:num],
						err:  err,
					}
				}
				n.recvCh <- &netPayload{
					data: b[:num],
					err:  err,
				}
				continue
			}
			return
		}
	}()

	for {
		select {
		case payload := <-n.recvCh:
			if payload.err == nil {
				p, err := packet.Unmarshal(payload.data)
				fmt.Printf("p: %v, err: %s", p, err)
			}
		case <-n.shutdownCh:
			return
		}
	}
}
