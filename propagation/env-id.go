package propagation

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

type envIDKeyType int

const envIDKey envIDKeyType = iota

const EnvIDHeader = "env-id"

type EnvID struct{}

var _ propagation.TextMapPropagator = EnvID{}

func (e EnvID) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	switch v := ctx.Value(envIDKey).(type) {
	case string:
		if v != "" {
			carrier.Set(EnvIDHeader, v)
		}
	}
}

func (e EnvID) Extract(parent context.Context, carrier propagation.TextMapCarrier) context.Context {
	v := carrier.Get(EnvIDHeader)
	if v == "" {
		return parent
	}
	return context.WithValue(parent, envIDKey, v)
}

func (e EnvID) Fields() []string {
	return []string{EnvIDHeader}
}
