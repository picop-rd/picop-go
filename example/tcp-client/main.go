package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/hiroyaonoe/bcop-go/protocol/header"
	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
	"go.opentelemetry.io/otel/baggage"
)

func main() {
	target, err := net.ResolveTCPAddr("tcp", "localhost:8080")
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

	_, err = io.Copy(bconn, os.Stdin)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}

}

func TestBaggage() baggage.Baggage {
	m1p1, _ := baggage.NewKeyProperty("p1Key")
	m1p2, _ := baggage.NewKeyValueProperty("p2Key", "p2Value")
	m1, _ := baggage.NewMember("m1Key", "m1Value", m1p1, m1p2)
	m2, _ := baggage.NewMember("m2Key", "m2Value")
	b, _ := baggage.New(m1, m2)
	return b
}
