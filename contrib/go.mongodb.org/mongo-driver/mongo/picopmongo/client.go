package picopmongo

import (
	"context"

	"github.com/picop-rd/picop-go/propagation"

	// Please use github.com/picop-rd/mongo-go-driver instead of go.mongodb.org/mongo-driver by replacing in go.mod
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	opts *options.ClientOptions
}

func New(opts *options.ClientOptions) *Client {
	if opts == nil {
		opts = options.Client()
	}
	opts = opts.
		SetDialer(DialContext(propagation.EnvID{})).
		SetMaxConnIdleTime(1).SetMaxPoolSize(1) // Requests with different headers must not use the same connection.
	return &Client{opts: opts}
}

func (c *Client) Connect(ctx context.Context) (*mongo.Client, error) {
	return mongo.Connect(ctx, c.opts)
}
