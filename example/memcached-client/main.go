package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
	picopmc "github.com/picop-rd/picop-go/contrib/github.com/bradfitz/gomemcache/picopgomemcache"
	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
)

func main() {
	// Prepare propagated context
	h := header.NewV1()
	h.Set(propagation.EnvIDHeader, "aaaaa")
	ctx := propagation.EnvID{}.Extract(context.Background(), propagation.NewPiCoPCarrier(h))

	uri := "localhost:11211"
	pc := picopmc.New(uri)
	defer pc.Close()

	mc := pc.Connect(ctx)
	err := mc.Set(ctx, &memcache.Item{
		Key:   "picop-example",
		Value: []byte("example"),
	})
	if err != nil {
		log.Fatal(err)
	}

	it, err := mc.Get(ctx, "picop-example")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(it.Value))
}
