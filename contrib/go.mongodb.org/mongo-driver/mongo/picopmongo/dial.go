package picopmongo

import (
	"context"
	"net"

	picopprop "github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	picopnet "github.com/picop-rd/picop-go/protocol/net"
	otelprop "go.opentelemetry.io/otel/propagation"

	// Please use github.com/picop-rd/mongo-go-driver instead of go.mongodb.org/mongo-driver by replacing in go.mod
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dialContext struct {
	options.ContextDialer
	propagator otelprop.TextMapPropagator
}

func DialContext(propagator otelprop.TextMapPropagator) dialContext {
	return dialContext{
		propagator: propagator,
	}
}

// DialContext is called by connection.connect
// https://github.com/mongodb/mongo-go-driver/blob/b8004e68a007b45b68a4cf39a03d98026b2230c9/x/mongo/driver/topology/connection.go#L196
func (d dialContext) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	nd := net.Dialer{}
	conn, err := nd.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}

	h := header.NewV1()
	d.propagator.Inject(ctx, picopprop.NewPiCoPCarrier(h))

	bconn := picopnet.SenderConn(conn, h)
	err = bconn.WriteHeader()
	if err != nil {
		return nil, err
	}
	return bconn, nil
}
