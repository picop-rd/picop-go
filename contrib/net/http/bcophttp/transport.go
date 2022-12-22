package bcophttp

import (
	"context"
	"net"
	"net/http"

	bcopprop "github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
	otelprop "go.opentelemetry.io/otel/propagation"
)

type Transport struct {
	*http.Transport
	propagator otelprop.TextMapPropagator
}

func NewTransport(base *http.Transport, propagator otelprop.TextMapPropagator) Transport {
	if base == nil {
		base = http.DefaultTransport.(*http.Transport)
	}

	t := Transport{
		Transport:  base.Clone(),
		propagator: propagator,
	}

	t.DialContext = wrapDialContext(base.DialContext, propagator)
	t.DialTLSContext = wrapDialContext(base.DialTLSContext, propagator)

	// 異なるヘッダのリクエスト同士が同じコネクション使ってはいけない
	t.DisableKeepAlives = true
	return t
}

func wrapDialContext(dc func(ctx context.Context, network, addr string) (net.Conn, error), propagator otelprop.TextMapPropagator) func(ctx context.Context, network, addr string) (net.Conn, error) {
	if dc == nil {
		return nil
	}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := dc(ctx, network, addr)
		if err != nil {
			return nil, err
		}

		h := header.NewV1("")
		propagator.Inject(ctx, bcopprop.NewBCoPCarrier(h))

		bconn := bcopnet.SenderConn(conn, h)
		return bconn, err
	}
}
