package picophttp

import (
	"context"
	"net"

	picopnet "github.com/picop-rd/picop-go/protocol/net"
)

type piCoPHeaderContextKeyType int

const PiCoPHeaderContextKey piCoPHeaderContextKeyType = iota

func ConnContext(ctx context.Context, c net.Conn) context.Context {
	if bc, ok := c.(*picopnet.Conn); ok {
		ctx = context.WithValue(ctx, PiCoPHeaderContextKey, bc.Header)
	}
	return ctx
}
