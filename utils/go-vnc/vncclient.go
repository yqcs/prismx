// VNC client implementation.

package vnc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"

	"golang.org/x/net/context"
	"prismx_cli/utils/go-vnc/go/metrics"
)

// Connect negotiates a connection to a VNC server.
func Connect(ctx context.Context, c net.Conn, cfg *ClientConfig) (*ClientConn, error) {
	conn := NewClientConn(c, cfg)

	if err := conn.processContext(ctx); err != nil {
		return nil, err
	}

	if err := conn.protocolVersionHandshake(ctx); err != nil {
		conn.Close()
		return nil, err
	}
	if err := conn.securityHandshake(); err != nil {
		conn.Close()
		return nil, err
	}
	if err := conn.securityResultHandshake(); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

// A ClientConfig structure is used to configure a ClientConn. After
// one has been passed to initialize a connection, it must not be modified.
type ClientConfig struct {
	secType uint8 // The negotiated security type.

	// A slice of ClientAuth methods. Only the first instance that is
	// suitable by the server will be used to authenticate.
	Auth []ClientAuth

	// Password for servers that require authentication.
	Password string

	// Exclusive determines whether the connection is shared with other
	// clients. If true, then all other clients connected will be
	// disconnected when a connection is established to the VNC server.
	Exclusive bool
}

// NewClientConfig returns a populated ClientConfig.
func NewClientConfig(p string) *ClientConfig {
	return &ClientConfig{
		Auth: []ClientAuth{
			&ClientAuthNone{},
			&ClientAuthVNC{p},
		},
		Password: p,
	}
}

// The ClientConn type holds client connection information.
type ClientConn struct {
	c               net.Conn
	config          *ClientConfig
	protocolVersion string

	// Name associated with the desktop, sent from the server.
	desktopName string

	// Height of the frame buffer in pixels, sent from the server.
	fbHeight uint16

	// Width of the frame buffer in pixels, sent from the server.
	fbWidth uint16

	// Track metrics on system performance.
	metrics map[string]metrics.Metric
}

func NewClientConn(c net.Conn, cfg *ClientConfig) *ClientConn {
	return &ClientConn{
		c:      c,
		config: cfg,
		metrics: map[string]metrics.Metric{
			"bytes-received": &metrics.Gauge{},
			"bytes-sent":     &metrics.Gauge{},
		},
	}
}

// Close a connection to a VNC server.
func (c *ClientConn) Close() error {
	return c.c.Close()
}

// DesktopName returns the server provided desktop name.
func (c *ClientConn) DesktopName() string {
	return c.desktopName
}

// setDesktopName stores the server provided desktop name.
func (c *ClientConn) setDesktopName(name string) {
	c.desktopName = name
}

// FramebufferHeight returns the server provided framebuffer height.
func (c *ClientConn) FramebufferHeight() uint16 {
	return c.fbHeight
}

// setFramebufferHeight stores the server provided framebuffer height.
func (c *ClientConn) setFramebufferHeight(height uint16) {
	c.fbHeight = height
}

// FramebufferWidth returns the server provided framebuffer width.
func (c *ClientConn) FramebufferWidth() uint16 {
	return c.fbWidth
}

// setFramebufferWidth stores the server provided framebuffer width.
func (c *ClientConn) setFramebufferWidth(width uint16) {
	c.fbWidth = width
}

// receive a packet from the network.
func (c *ClientConn) receive(data interface{}) error {
	if err := binary.Read(c.c, binary.BigEndian, data); err != nil {
		return err
	}
	c.metrics["bytes-received"].Adjust(int64(binary.Size(data)))
	return nil
}

// receiveN receives N packets from the network.
func (c *ClientConn) receiveN(data interface{}, n int) error {
	if n == 0 {
		return nil
	}

	switch data.(type) {
	case *[]uint8:
		var v uint8
		for i := 0; i < n; i++ {
			if err := binary.Read(c.c, binary.BigEndian, &v); err != nil {
				return err
			}
			slice := data.(*[]uint8)
			*slice = append(*slice, v)
		}
	case *[]int32:
		var v int32
		for i := 0; i < n; i++ {
			if err := binary.Read(c.c, binary.BigEndian, &v); err != nil {
				return err
			}
			slice := data.(*[]int32)
			*slice = append(*slice, v)
		}
	case *bytes.Buffer:
		var v byte
		for i := 0; i < n; i++ {
			if err := binary.Read(c.c, binary.BigEndian, &v); err != nil {
				return err
			}
			buf := data.(*bytes.Buffer)
			buf.WriteByte(v)
		}
	default:
		return errors.New(fmt.Sprintf("unrecognized data type %v", reflect.TypeOf(data)))
	}
	c.metrics["bytes-received"].Adjust(int64(binary.Size(data)))
	return nil
}

// send a packet to the network.
func (c *ClientConn) send(data interface{}) error {
	if err := binary.Write(c.c, binary.BigEndian, data); err != nil {
		return err
	}
	c.metrics["bytes-sent"].Adjust(int64(binary.Size(data)))
	return nil
}

func (c *ClientConn) processContext(ctx context.Context) error {
	if mpv := ctx.Value("vnc_max_proto_version"); mpv != nil && mpv != "" {
		log.Printf("vnc_max_proto_version: %v", mpv)
		vers := []string{"3.3", "3.8"}
		valid := false
		for _, v := range vers {
			if mpv == v {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("Invalid max protocol version %v; supported versions are %v", mpv, vers)
		}
	}

	return nil
}

func (c *ClientConn) DebugMetrics() {
	log.Println("Metrics:")
	for name, metric := range c.metrics {
		log.Printf("  %v: %v", name, metric.Value())
	}
}
