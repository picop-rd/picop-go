package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	bcopprop "github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	otelprop "go.opentelemetry.io/otel/propagation"
)

func main() {
	// 伝播されたContextを用意
	h := header.NewV1("key1=value1")
	ctx := otelprop.Baggage{}.Extract(context.Background(), bcopprop.NewBCoPCarrier(h))

	client := &http.Client{
		Transport: bcophttp.NewTransport(nil, otelprop.Baggage{}),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
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
