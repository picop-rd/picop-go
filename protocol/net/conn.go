package net

import (
	"bufio"
	"net"
	"sync"

	"github.com/hiroyaonoe/bcop-go/protocol/header"
)

type connType int

const (
	typeSender connType = iota
	typeReceiver
)

type Conn struct {
	net.Conn
	connType  connType
	bufReader *bufio.Reader
	Header    *header.Header
	once      sync.Once
	err       error
}

var _ net.Conn = &Conn{}

func SenderConn(c net.Conn, h *header.Header) *Conn {
	return &Conn{
		Conn:      c,
		connType:  typeSender,
		bufReader: bufio.NewReader(c),
		Header:    h,
	}
}

func ReceiverConn(c net.Conn) *Conn {
	return &Conn{
		Conn:      c,
		connType:  typeReceiver,
		bufReader: bufio.NewReader(c),
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	c.readHeader()
	if c.err != nil {
		return 0, c.err
	}

	return c.bufReader.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.writeHeader()
	if c.err != nil {
		return 0, c.err
	}

	return c.Conn.Write(b)
}

func (c *Conn) ReadHeader() (*header.Header, error) {
	c.readHeader()
	if c.err != nil {
		return nil, c.err
	}
	return c.Header, nil
}

func (c *Conn) readHeader() {
	if c.connType == typeReceiver {
		c.once.Do(func() {
			c.Header, c.err = header.Parse(c.bufReader)
		})
	}
}

func (c *Conn) writeHeader() {
	if c.connType == typeSender {
		c.once.Do(func() {
			_, c.err = c.Header.WriteTo(c.Conn)
		})
	}
}
