package picopmongo

import (
	"context"
	"sync"

	"github.com/picop-rd/picop-go/propagation"
	picopprop "github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"

	// Please use github.com/picop-rd/mongo-go-driver instead of go.mongodb.org/mongo-driver by replacing in go.mod
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	opts        *options.ClientOptions
	pool        *sync.Map
	PoolByEnvID bool
}

func New(opts *options.ClientOptions) *Client {
	if opts == nil {
		opts = options.Client()
	}
	opts = opts.
		SetDialer(DialContext(propagation.EnvID{}))
	return &Client{
		opts:        opts,
		pool:        &sync.Map{},
		PoolByEnvID: true,
	}
}

func (c *Client) DisablePoolByEnvID() {
	c.PoolByEnvID = false
}

func (c *Client) Connect(ctx context.Context) (*mongo.Client, error) {
	if !c.PoolByEnvID {
		return mongo.Connect(ctx, c.opts)
	}

	h := header.NewV1()
	propagation.EnvID{}.Inject(ctx, picopprop.NewPiCoPCarrier(h))
	envID := h.Get(propagation.EnvIDHeader)

	if client, ok := c.pool.Load(envID); ok {
		return client.(*mongo.Client), nil
	}
	nc, err := mongo.Connect(ctx, c.opts)
	if err != nil {
		return nil, err
	}
	c.pool.Store(envID, nc)
	return nc, nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	if !c.PoolByEnvID {
		return nil
	}

	var err error
	c.pool.Range(func(key, value interface{}) bool {
		err = value.(*mongo.Client).Disconnect(ctx)
		return err == nil
	})
	return err
}
