package main

import (
	"context"

	"github.com/picop-rd/picop-go/contrib/google.golang.org/grpc/picopgrpc"
	"github.com/picop-rd/picop-go/example/grpc-server/proto"
	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Prepare propagated context
	h := header.NewV1()
	h.Set(propagation.EnvIDHeader, "aaaaa")
	ctx := propagation.EnvID{}.Extract(context.Background(), propagation.NewPiCoPCarrier(h))
	ctx = context.WithValue(ctx, "picop", "propagated")
	url := "localhost:8080"
	pc := picopgrpc.New(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer pc.Close()

	cc, err := pc.Connect(ctx)
	if err != nil {
		panic(err)
	}

	ec := proto.NewExampleClient(cc)
	res, err := ec.Get(ctx, &proto.Request{
		Id: "picop",
	})
	if err != nil {
		panic(err)
	}
	println(res.Id)
}
