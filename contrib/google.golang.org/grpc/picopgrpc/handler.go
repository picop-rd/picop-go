package picopgrpc

import (
	"context"

	"github.com/picop-rd/picop-go/propagation"
	picopprop "github.com/picop-rd/picop-go/propagation"
	picopnet "github.com/picop-rd/picop-go/protocol/net"
	otelprop "go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(propagator otelprop.TextMapPropagator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		c := grpc.GetConnection(ctx)
		bc, ok := c.(*picopnet.Conn)
		if !ok {
			return handler(ctx, req)
		}
		hd := bc.Header
		ctx = propagation.EnvID{}.Extract(ctx, picopprop.NewPiCoPCarrier(hd))
		return handler(ctx, req)
	}
}
