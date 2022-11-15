package protocol

import (
	"bufio"
	"net"

	"github.com/hiroyaonoe/bopp-go/header"
)

type Conn struct {
	net.Conn
	header    *header.Header
	bufReader *bufio.Reader
}

var _ net.Conn = &Conn{}

func NewConn(c net.Conn, h *header.Header) *Conn {
	return &Conn{
		Conn:      c,
		header:    h,
		bufReader: bufio.NewReader(c),
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	err = c.readHeader()
	if err != nil {
		return 0, err
	}

	return c.bufReader.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	_, err = c.header.WriteTo(c.Conn)
	if err != nil {
		return 0, err
	}
	n, err = c.Write(b)
	return n, err
}

func (c *Conn) readHeader() error {
	h, err := header.Read(c.bufReader)
	if err != nil {
		return err
	}
	c.header = h
	return nil
}
