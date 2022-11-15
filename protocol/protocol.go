package protocol

import (
	"net"

	"github.com/hiroyaonoe/bopp-go/header"
)

type Conn struct {
	net.Conn
	header header.Header
}

var _ net.Conn = &Conn{}
