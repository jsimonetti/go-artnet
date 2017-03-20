package packet

import (
	"net"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ListenRawUDP4(ip net.Addr) (net.PacketConn, error) {
	srcIP := ip.(*net.UDPAddr).IP
	addr := &net.IPAddr{
		IP: srcIP.To4(),
	}

	c, err := net.ListenIP("ip4:udp", addr)
	if err != nil {
		return nil, err
	}

	conn := &conn{
		conn:  c,
		src:   ip,
		srcIP: srcIP,
		buf:   gopacket.NewSerializeBuffer(),
		opts: gopacket.SerializeOptions{
			ComputeChecksums: true,
			FixLengths:       true,
		},
	}

	return conn, err
}

type conn struct {
	conn *net.IPConn

	src   net.Addr
	srcIP net.IP

	buf  gopacket.SerializeBuffer
	opts gopacket.SerializeOptions
}

func (c *conn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	dst := addr.(*net.IPAddr).IP

	// build the ip layer
	ipLayer := layers.IPv4{
		SrcIP:    c.srcIP,
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

func (c *conn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	var nread int
	buf := make([]byte, 1500)

	nread, addr, err = c.conn.ReadFromIP(buf)
	if err != nil {
		return
	}

	//var ip4 layers.IPv4
	var udp layers.UDP
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeUDP, &udp)
	decoded := []gopacket.LayerType{}

	err = parser.DecodeLayers(buf[:nread], &decoded)

	spew.Dump(udp.Payload)

	if err != nil {
		return
	}
	n = copy(b, udp.Payload)

	return
}

func (c *conn) Close() error {
	return c.conn.Close()
}

func (c *conn) LocalAddr() net.Addr {
	return c.src
}

func (c *conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
