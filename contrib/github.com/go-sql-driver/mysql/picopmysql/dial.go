package picopmysql

import (
	"context"
	"net"

	"github.com/go-sql-driver/mysql"
	picopprop "github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	picopnet "github.com/picop-rd/picop-go/protocol/net"
	otelprop "go.opentelemetry.io/otel/propagation"
)

//TODO: propagatorが固定できる場合はinit()でRegisterDialContext("tcp", propagator)する

func RegisterDialContext(netP string, propagator otelprop.TextMapPropagator) {
	mysql.RegisterDialContext(netP, DialContext(netP, propagator))
}

func DialContext(netP string, propagator otelprop.TextMapPropagator) mysql.DialContextFunc {
	// connector.Connectで呼び出される
	// https://github.com/go-sql-driver/mysql/blob/4591e42e65cf483147a7c7a4f4cfeac81b21c917/connector.go#L37
	return func(ctx context.Context, addr string) (net.Conn, error) {
		nd := net.Dialer{}
		conn, err := nd.DialContext(ctx, netP, addr)
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
