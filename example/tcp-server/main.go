package main

import (
	"io"
	"log"
	"net"
	"os"

	bcopnet "github.com/hiroyaonoe/bcop-go/net"
)

func main() {
	addr := "localhost:8080"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	conn, err := ln.Accept()
	defer conn.Close()

	bconn := bcopnet.NewConn(conn, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(os.Stdout, bconn)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatalf(err.Error())
	}
}
