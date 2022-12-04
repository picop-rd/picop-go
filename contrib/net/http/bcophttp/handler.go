package bcophttp

import (
	"net/http"

	bcopprop "github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	otelprop "go.opentelemetry.io/otel/propagation"
)

type Handler struct {
	http.Handler
	propagator otelprop.TextMapPropagator
}

var _ http.Handler = Handler{}

func NewHandler(hl http.Handler, prop otelprop.TextMapPropagator) Handler {
	return Handler{
		Handler:    hl,
		propagator: prop,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hd := ctx.Value(BCoPHeaderContextKey).(*header.Header)
	ctx = h.propagator.Extract(ctx, bcopprop.NewBCoPCarrier(hd))
	nr := r.Clone(ctx)

	h.Handler.ServeHTTP(w, nr)
}
