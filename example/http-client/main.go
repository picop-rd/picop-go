package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	bcopprop "github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	"go.opentelemetry.io/otel/baggage"
	otelprop "go.opentelemetry.io/otel/propagation"
)

func main() {
	port := flag.String("port", "8080", "request port")
	flag.Parse()
	// 伝播されたContextを用意
	bag := TestBaggage()
	h := header.NewV1(bag.String())
	ctx := otelprop.Baggage{}.Extract(context.Background(), bcopprop.NewBCoPCarrier(h))

	client := &http.Client{
		Transport: bcophttp.NewTransport(nil, otelprop.Baggage{}),
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

func TestBaggage() baggage.Baggage {
	m1p1, _ := baggage.NewKeyProperty("p1Key")
	m1p2, _ := baggage.NewKeyValueProperty("p2Key", "p2Value")
	m1, _ := baggage.NewMember("m1Key", "m1Value", m1p1, m1p2)
	m2, _ := baggage.NewMember("m2Key", "m2Value")
	m3, _ := baggage.NewMember("env-id", "aaaaa")
	b, _ := baggage.New(m1, m2, m3)
	return b
}
