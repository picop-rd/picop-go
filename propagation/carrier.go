package propagation

import (
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	"go.opentelemetry.io/otel/propagation"
)

const baggageKey = "baggage"

// BCoPCarrier adapts BCoP to satisfy the OpenTelemetry TextMapCarrier interface.
// Warning: The key must be "baggage". In the future, this will be an HTTP Header compliant key-value store.
// (See https://opentelemetry.io/docs/reference/specification/context/api-propagators/#textmap-propagator )
type BCoPCarrier struct {
	*header.Header
}

var _ propagation.TextMapCarrier = BCoPCarrier{}

func NewBCoPCarrier(h *header.Header) BCoPCarrier {
	return BCoPCarrier{h}
}

func (bc BCoPCarrier) Get(key string) string {
	if key == baggageKey {
		return bc.Header.Get()
	}
	return ""
}

func (bc BCoPCarrier) Set(key, value string) {
	if key == baggageKey {
		bc.Header.Set(value)
		return
	}
	return
}

func (bc BCoPCarrier) Keys() []string {
	return []string{baggageKey}
}
