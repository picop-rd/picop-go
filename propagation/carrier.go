package propagation

import (
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	"go.opentelemetry.io/otel/propagation"
)

// BCoPCarrier adapts BCoP to satisfy the OpenTelemetry TextMapCarrier interface.
type BCoPCarrier struct {
	*header.Header
}

var _ propagation.TextMapCarrier = BCoPCarrier{}

func NewBCoPCarrier(h *header.Header) BCoPCarrier {
	return BCoPCarrier{h}
}

func (bc BCoPCarrier) Get(key string) string {
	return bc.Header.Get(key)
}

func (bc BCoPCarrier) Set(key, value string) {
	bc.Header.Set(key, value)
	return
}

func (bc BCoPCarrier) Keys() []string {
	return bc.Header.Keys()
}
