package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/picop-rd/picop-go/contrib/google.golang.org/grpc/picopgrpc"
	"github.com/picop-rd/picop-go/example/grpc-server/proto"
	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	picopnet "github.com/picop-rd/picop-go/protocol/net"
	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "8080", "listen port")
	flag.Parse()

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			picopgrpc.UnaryServerInterceptor(propagation.EnvID{}),
		),
	}

	srv := grpc.NewServer(opts...)
	s := &exampleServer{}
	proto.RegisterExampleServer(srv, s)

	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	bln := picopnet.NewListener(ln)

	log.Fatal(srv.Serve(bln))
}

type exampleServer struct {
	proto.ExampleServer
}

func (s *exampleServer) Get(ctx context.Context, req *proto.Request) (*proto.Result, error) {
	// Confirm propagated context
	h := header.NewV1()
	propagation.EnvID{}.Inject(ctx, propagation.NewPiCoPCarrier(h))
	fmt.Printf("PiCoP Header Accepted: %s\n", h)

	return &proto.Result{
		Id: req.Id,
	}, nil
}
