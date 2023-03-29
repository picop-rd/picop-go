package picophttp

import (
	"net/http"

	picopprop "github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
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
	hd := ctx.Value(PiCoPHeaderContextKey).(*header.Header)
	ctx = h.propagator.Extract(ctx, picopprop.NewPiCoPCarrier(hd))
	nr := r.Clone(ctx)

	h.Handler.ServeHTTP(w, nr)
}
