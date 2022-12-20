package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	bcopprop "github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
	otelprop "go.opentelemetry.io/otel/propagation"
)

func main() {
	port := flag.String("port", "8080", "listen port")
	flag.Parse()

	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:        ":" + *port,
		Handler:     bcophttp.NewHandler(http.DefaultServeMux, otelprop.Baggage{}),
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
	h := header.NewV1("")
	otelprop.Baggage{}.Inject(r.Context(), bcopprop.NewBCoPCarrier(h))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("BCoP Header Accepted: %s\n", h))
}
