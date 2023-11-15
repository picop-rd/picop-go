package picopgrpc

import (
	"context"
	"net"
	"time"

	picopprop "github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	picopnet "github.com/picop-rd/picop-go/protocol/net"
	otelprop "go.opentelemetry.io/otel/propagation"
)

// DialContext is called by dial
// https://github.com/grpc/grpc-go/blob/8645f95509d6c5d17a54621407f3ca717d4f8620/internal/transport/http2_client.go#L152
func DialContext(propagator otelprop.TextMapPropagator) func(ctx context.Context, address string) (net.Conn, error) {
	return func(ctx context.Context, address string) (net.Conn, error) {
		nd := (&net.Dialer{KeepAlive: time.Duration(-1)})
		conn, err := nd.DialContext(ctx, "tcp", address)
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
