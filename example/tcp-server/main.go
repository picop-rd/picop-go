package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	picopnet "github.com/picop-rd/picop-go/protocol/net"
	"golang.org/x/sync/errgroup"
)

func main() {
	port := flag.String("port", "8080", "listen port")
	flag.Parse()

	addr := "localhost:" + *port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	bln := picopnet.NewListener(ln)

	bconn, err := bln.AcceptWithPiCoPConn()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer bconn.Close()

	header, err := bconn.ReadHeader()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PiCoP Header Accepted: %s\n", header)

	var eg errgroup.Group

	eg.Go(func() error { _, err := io.Copy(bconn, os.Stdin); return err })
	eg.Go(func() error { _, err := io.Copy(os.Stdout, bconn); return err })

	err = eg.Wait()
	log.Fatal(err)
}
