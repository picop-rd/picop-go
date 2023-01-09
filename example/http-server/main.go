package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	"github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
)

func main() {
	port := flag.String("port", "8080", "listen port")
	flag.Parse()

	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:        ":" + *port,
		Handler:     bcophttp.NewHandler(http.DefaultServeMux, propagation.EnvID{}),
		ConnContext: bcophttp.ConnContext,
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	bln := bcopnet.NewListener(ln)

	log.Fatal(server.Serve(bln))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 伝播されたContextを確認
	h := header.NewV1()
	propagation.EnvID{}.Inject(r.Context(), propagation.NewBCoPCarrier(h))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("BCoP Header Accepted: %s\n", h))
}
