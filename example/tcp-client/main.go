package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/hiroyaonoe/bcop-go/protocol/header"
	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
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

	h := &header.Header{}
	bconn := bcopnet.NewConn(conn, h)

	_, err = io.Copy(bconn, os.Stdin)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}

}
