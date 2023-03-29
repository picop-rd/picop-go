package propagation

import (
	"github.com/picop-rd/picop-go/protocol/header"
	"go.opentelemetry.io/otel/propagation"
)

// PiCoPCarrier adapts PiCoP to satisfy the OpenTelemetry TextMapCarrier interface.
type PiCoPCarrier struct {
	*header.Header
}

var _ propagation.TextMapCarrier = PiCoPCarrier{}

func NewPiCoPCarrier(h *header.Header) PiCoPCarrier {
	return PiCoPCarrier{h}
}

func (bc PiCoPCarrier) Get(key string) string {
	return bc.Header.Get(key)
}

func (bc PiCoPCarrier) Set(key, value string) {
	bc.Header.Set(key, value)
	return
}

func (bc PiCoPCarrier) Keys() []string {
	return bc.Header.Keys()
}
