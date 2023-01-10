package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	"github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
)

func main() {
	port := flag.String("port", "8080", "request port")
	flag.Parse()
	// 伝播されたContextを用意
	h := header.NewV1()
	h.Set(propagation.EnvIDHeader, "aaaaa")
	ctx := propagation.EnvID{}.Extract(context.Background(), propagation.NewBCoPCarrier(h))

	client := &http.Client{
		Transport: bcophttp.NewTransport(nil, propagation.EnvID{}),
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
