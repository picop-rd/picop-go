package net

import (
	"net"
)

type Listener struct {
	net.Listener
}

var _ net.Listener = Listener{}

func NewListener(l net.Listener) Listener {
	return Listener{l}
}
func (l Listener) AcceptWithBCoPConn() (*Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	bconn := ReceiverConn(conn)
	bconn.readHeader()
	if bconn.err != nil {
		return nil, err
	}

	return bconn, nil
}

func (l Listener) Accept() (net.Conn, error) {
	return l.AcceptWithBCoPConn()
}
