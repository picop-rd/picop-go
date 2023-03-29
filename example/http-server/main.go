package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/picop-rd/picop-go/contrib/net/http/picophttp"
	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	picopnet "github.com/picop-rd/picop-go/protocol/net"
)

func main() {
	port := flag.String("port", "8080", "listen port")
	flag.Parse()

	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:        ":" + *port,
		Handler:     picophttp.NewHandler(http.DefaultServeMux, propagation.EnvID{}),
		ConnContext: picophttp.ConnContext,
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	bln := picopnet.NewListener(ln)

	log.Fatal(server.Serve(bln))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 伝播されたContextを確認
	h := header.NewV1()
	propagation.EnvID{}.Inject(r.Context(), propagation.NewPiCoPCarrier(h))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("PiCoP Header Accepted: %s\n", h))
}
