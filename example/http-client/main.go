package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/picop-rd/picop-go/contrib/net/http/picophttp"
	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
)

func main() {
	port := flag.String("port", "8080", "request port")
	flag.Parse()
	// Prepare propagated context
	h := header.NewV1()
	h.Set(propagation.EnvIDHeader, "aaaaa")
	ctx := propagation.EnvID{}.Extract(context.Background(), propagation.NewPiCoPCarrier(h))

	client := &http.Client{
		Transport: picophttp.NewTransport(nil, propagation.EnvID{}),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
