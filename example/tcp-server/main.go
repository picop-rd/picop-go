package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
	"golang.org/x/sync/errgroup"
)

func main() {
	addr := "localhost:8080"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	bln := bcopnet.NewListener(ln)

	bconn, err := bln.AcceptWithBCoPConn()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer bconn.Close()

	header, err := bconn.ReadHeader()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BCoP Header Accepted: %s\n", header)

	var eg errgroup.Group

	eg.Go(func() error { _, err := io.Copy(bconn, os.Stdin); return err })
	eg.Go(func() error { _, err := io.Copy(os.Stdout, bconn); return err })

	err = eg.Wait()
	log.Fatal(err)
}
