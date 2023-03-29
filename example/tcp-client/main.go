package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	picopnet "github.com/picop-rd/picop-go/protocol/net"
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

	h := header.NewV1()
	h.Set(propagation.EnvIDHeader, "aaaaa")
	bconn := picopnet.SenderConn(conn, h)

	var eg errgroup.Group

	eg.Go(func() error { _, err := io.Copy(bconn, os.Stdin); return err })
	eg.Go(func() error { _, err := io.Copy(os.Stdout, bconn); return err })

	err = eg.Wait()
	log.Fatal(err)
}
