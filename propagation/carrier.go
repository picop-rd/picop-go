package propagation

import (
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	"go.opentelemetry.io/otel/propagation"
)

const baggageKey = "baggage"

// BCoPCarrier adapts BCoP to satisfy the OpenTelemetry TextMapCarrier interface.
// Warning: The key must be "baggage". In the future, this will be an HTTP Header compliant key-value store.
// (See https://opentelemetry.io/docs/reference/specification/context/api-propagators/#textmap-propagator )
type BCoPCarrier header.Header

var _ propagation.TextMapCarrier = BCoPCarrier{}

func (bc BCoPCarrier) Get(key string) string {
	if key == baggageKey {
		return bc.Value
	}
	return ""
}

func (bc BCoPCarrier) Set(key, value string) {
	if key == baggageKey {
		bc.Value = value
		return
	}
	return
}

func (bc BCoPCarrier) Keys() []string {
	return []string{baggageKey}
}
