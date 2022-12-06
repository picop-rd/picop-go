package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
)

func main() {
	addr := "localhost:8080"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	bln := bcopnet.NewListener(ln)

	bconn, err := bln.AcceptWithBCoPConn()
	defer bconn.Close()

	header, err := bconn.ReadHeader()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BCoP Header Accepted: %s\n", header)

	_, err = io.Copy(os.Stdout, bconn)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatalf(err.Error())
	}
}
