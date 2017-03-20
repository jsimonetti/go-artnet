package packet

import (
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// ListenRawUDP4 will listen for ipv4 udp packets on ip
// It is a custom connection implementing the PacketConn interface
func ListenRawUDP4(ip net.Addr) (net.PacketConn, error) {
	addr := &net.IPAddr{
		IP: ip.(*net.UDPAddr).IP.To4(),
	}

	// start a listener on the specified address
	c, err := net.ListenIP("ip4:udp", addr)
	if err != nil {
		return nil, err
	}

	conn := &conn{
		conn: c,
		src:  ip,
		buf:  gopacket.NewSerializeBuffer(),
		opts: gopacket.SerializeOptions{
			ComputeChecksums: true,
			FixLengths:       true,
		},
	}

	return conn, err
}

// conn holds our connection
type conn struct {
	// the raw IP connection
	conn *net.IPConn

	// the ip address of the listener
	src net.Addr

	// used for creating packets
	buf  gopacket.SerializeBuffer
	opts gopacket.SerializeOptions
}

// ReadFrom reads a packet from the connection,
// copying the payload into b. It returns the number of
// bytes copied into b and the return address that
// was on the packet.
// ReadFrom can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetReadDeadline.
func (c *conn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	buf := make([]byte, 4096)

	_, addr, err = c.conn.ReadFromIP(buf)
	if err != nil {
		return
	}

	// decode the udp packet
	var udp layers.UDP
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeUDP, &udp)
	decoded := []gopacket.LayerType{}

	err = parser.DecodeLayers(buf, &decoded)
	// the "No decoder for layer type Payload" is expected as we don't use a decode for the UDP payload
	if err != nil && err.Error() != "No decoder for layer type Payload" {
		return
	}

	// copy the contents back
	n = copy(b, udp.Payload)

	return
}

// WriteTo writes a packet with payload b to addr.
// WriteTo can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetWriteDeadline.
// On packet-oriented connections, write timeouts are rare.
func (c *conn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	dst := addr.(*net.IPAddr).IP

	// build the ip layer
	ipLayer := layers.IPv4{
		SrcIP:    c.src.(*net.UDPAddr).IP,
		DstIP:    dst,
		Protocol: layers.IPProtocolUDP,
	}

	// build the udp layer
	udpLayer := layers.UDP{
		SrcPort: ArtNetPort,
		DstPort: ArtNetPort,
	}
	udpLayer.SetNetworkLayerForChecksum(&ipLayer)

	// serialize the layer into a buffer
	err = gopacket.SerializeLayers(c.buf, c.opts, &udpLayer, gopacket.Payload(b))
	if err != nil {
		return
	}

	to := &net.IPAddr{
		IP: dst,
	}

	// write the buffer into the conn
	n, err = c.conn.WriteToIP(c.buf.Bytes(), to)
	n -= 8
	return
}

// Close closes the connection.
// Any blocked ReadFrom or WriteTo operations will be unblocked and return errors.
func (c *conn) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *conn) LocalAddr() net.Addr {
	return c.src
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to ReadFrom or
// WriteTo. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful ReadFrom or WriteTo calls.
//
// A zero value for t means I/O operations will not time out.
func (c *conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future ReadFrom calls
// and any currently-blocked ReadFrom call.
// A zero value for t means ReadFrom will not time out.
func (c *conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future WriteTo calls
// and any currently-blocked WriteTo call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means WriteTo will not time out.
func (c *conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
