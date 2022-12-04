package bcophttp

import (
	"context"
	"net"

	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
)

type bCoPHeaderContextKeyType int

const BCoPHeaderContextKey bCoPHeaderContextKeyType = iota

func ConnContext(ctx context.Context, c net.Conn) context.Context {
	if bc, ok := c.(*bcopnet.Conn); ok {
		ctx = context.WithValue(ctx, BCoPHeaderContextKey, bc.Header)
	}
	return ctx
}
