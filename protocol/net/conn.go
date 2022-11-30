package protocol

import (
	"bufio"
	"net"
	"sync"

	"github.com/hiroyaonoe/bcop-go/protocol/header"
)

type Conn struct {
	net.Conn
	bufReader *bufio.Reader
	rHeader   *header.Header
	rOnce     sync.Once
	rErr      error
	wHeader   *header.Header
	wOnce     sync.Once
	wErr      error
}

var _ net.Conn = &Conn{}

func NewConn(c net.Conn, h *header.Header) *Conn {
	return &Conn{
		Conn:      c,
		bufReader: bufio.NewReader(c),
		wHeader:   h,
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	c.readHeader()
	if c.rErr != nil {
		return 0, c.rErr
	}

	return c.bufReader.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.writeHeader()
	if c.wErr != nil {
		return 0, c.wErr
	}

	return c.Conn.Write(b)
}

func (c *Conn) ReadHeader() (*header.Header, error) {
	c.readHeader()
	if c.rErr != nil {
		return nil, c.rErr
	}
	return c.rHeader, nil
}

func (c *Conn) readHeader() {
	c.rOnce.Do(func() {
		c.rHeader, c.rErr = header.Parse(c.bufReader)
	})
}

func (c *Conn) writeHeader() {
	c.wOnce.Do(func() {
		_, c.wErr = c.wHeader.WriteTo(c.Conn)
	})
}
