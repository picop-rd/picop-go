package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/hiroyaonoe/bcop-go/protocol/header"
	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
	"go.opentelemetry.io/otel/baggage"
	"golang.org/x/sync/errgroup"
)

func main() {
	port := flag.String("port", "8080", "dial port")
	flag.Parse()
	target, err := net.ResolveTCPAddr("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, target)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	bag := TestBaggage()
	h := header.NewV1(bag.String())
	bconn := bcopnet.SenderConn(conn, h)

	var eg errgroup.Group

	eg.Go(func() error { _, err := io.Copy(bconn, os.Stdin); return err })
	eg.Go(func() error { _, err := io.Copy(os.Stdout, bconn); return err })

	err = eg.Wait()
	log.Fatal(err)
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
