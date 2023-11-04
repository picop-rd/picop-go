package picopgomemcache

import (
	"context"
	"net"

	picopprop "github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	picopnet "github.com/picop-rd/picop-go/protocol/net"
	otelprop "go.opentelemetry.io/otel/propagation"
)

// DialContext is called by Client.dial
// https://github.com/bradfitz/gomemcache/blob/24af94b0387418c51cc45a2e1fe6d4d1bef8a0fd/memcache/memcache.go#L280
func DialContext(propagator otelprop.TextMapPropagator) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		nd := net.Dialer{}
		conn, err := nd.DialContext(ctx, network, address)
		if err != nil {
			return nil, err
		}

		h := header.NewV1()
		propagator.Inject(ctx, picopprop.NewPiCoPCarrier(h))

		bconn := picopnet.SenderConn(conn, h)
		err = bconn.WriteHeader()
		if err != nil {
			return nil, err
		}
		return bconn, nil
	}
}
