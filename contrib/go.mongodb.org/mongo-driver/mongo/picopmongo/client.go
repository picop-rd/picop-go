package picopmongo

import (
	"context"

	"github.com/picop-rd/picop-go/propagation"
	otelprop "go.opentelemetry.io/otel/propagation"

	// Please use github.com/picop-rd/mongo-go-driver instead of go.mongodb.org/mongo-driver by replacing in go.mod
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	opts *options.ClientOptions
}

func New(opts *options.ClientOptions, propagator otelprop.TextMapPropagator) *Client {
	if opts == nil {
		opts = options.Client()
	}
	opts = opts.
		SetDialer(DialContext(propagation.EnvID{}))
	return &Client{opts: opts}
}

func (c *Client) Connect(ctx context.Context) (*mongo.Client, error) {
	return mongo.Connect(ctx, c.opts)
}
