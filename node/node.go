package node

import (
	"fmt"
	"net"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
)

// Node is the information known about a node
type Node struct {
	// Config holds the configuration of this node
	Config NodeConfig

	// conn is the UDP connection this node will listen on
	conn   net.Conn
	sendCh chan netPayload
	recvCh chan netPayload

	// shutdownCh will be closed on shutdown of the node
	shutdownCh chan struct{}

	// pollCh will receive ArtPoll packets
	pollCh chan packet.ArtPollPacket

	// Controller is a config of a controller should this node by under it's controller
	Controller NodeConfig
}

type netPayload struct {
	err  error
	data []byte
}

// New return a Node
func New() Node {
	return Node{
		Config:     NodeConfig{},
		sendCh:     make(chan netPayload),
		recvCh:     make(chan netPayload),
		pollCh:     make(chan packet.ArtPollPacket),
		shutdownCh: make(chan struct{}),
	}
}

// Close will stop all running routines and close this controller
func (n *Node) Close() {
	close(n.shutdownCh)
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
		case <-n.shutdownCh:
			return
		}
	}
}

func (n *Node) recvLoop() {
	for {
		select {
		case <-n.shutdownCh:
			return
		}
	}
}

func (n *Node) Write(b []byte) (num int, err error) {
	return 0, nil
}

func (n *Node) Read(b []byte) (num int, err error) {
	return 0, nil
}
