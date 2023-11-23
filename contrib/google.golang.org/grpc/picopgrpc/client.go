package picopgrpc

import (
	"context"
	"sync"

	"github.com/picop-rd/picop-go/propagation"
	picopprop "github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Client struct {
	opts        []grpc.DialOption
	target      string
	pool        *sync.Map
	PoolByEnvID bool
}

func New(target string, opts ...grpc.DialOption) *Client {
	return &Client{
		opts:        append(opts, grpc.WithContextDialer(DialContext(propagation.EnvID{}))),
		target:      target,
		pool:        &sync.Map{},
		PoolByEnvID: true,
	}
}

func (c *Client) DisablePoolByEnvID() {
	c.PoolByEnvID = false
}

func (c *Client) Connect(ctx context.Context) (*grpc.ClientConn, error) {
	if !c.PoolByEnvID {
		return grpc.DialContext(ctx, c.target, c.opts...)
	}

	h := header.NewV1()
	propagation.EnvID{}.Inject(ctx, picopprop.NewPiCoPCarrier(h))
	envID := h.Get(propagation.EnvIDHeader)

	if client, ok := c.pool.Load(envID); ok {
		if client.(*grpc.ClientConn).GetState() != connectivity.Shutdown {
			return client.(*grpc.ClientConn), nil
		}
	}
	cc, err := grpc.DialContext(ctx, c.target, c.opts...)
	if err != nil {
		return nil, err
	}
	c.pool.Store(envID, cc)
	return cc, err
}

func (c *Client) Close() error {
	if !c.PoolByEnvID {
		return nil
	}

	var err error
	c.pool.Range(func(key, value interface{}) bool {
		err = value.(*grpc.ClientConn).Close()
		return err == nil
	})
	return err
}
