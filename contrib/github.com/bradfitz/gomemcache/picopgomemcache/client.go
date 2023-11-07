package picopgomemcache

import (
	"context"
	"sync"

	// Please use github.com/picop-rd/gomemcache instead of github.com/bradfitz/gomemcache/memcache by replacing in go.mod
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/picop-rd/picop-go/propagation"
	picopprop "github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
)

type Client struct {
	server      []string
	pool        *sync.Map
	PoolByEnvID bool
}

func New(server ...string) *Client {
	return &Client{
		server: server,
	}
}

func (c *Client) DisablePoolByEnvID() {
	c.PoolByEnvID = false
}

func (c *Client) Connect(ctx context.Context) *memcache.Client {
	if !c.PoolByEnvID {
		mc := memcache.New(c.server...)
		mc.DialContext = DialContext(propagation.EnvID{})
		return mc
	}

	h := header.NewV1()
	propagation.EnvID{}.Inject(ctx, picopprop.NewPiCoPCarrier(h))
	envID := h.Get(propagation.EnvIDHeader)

	if client, ok := c.pool.Load(envID); ok {
		return client.(*memcache.Client)
	}
	mc := memcache.New(c.server...)
	mc.DialContext = DialContext(propagation.EnvID{})
	c.pool.Store(envID, mc)
	return mc
}

func (c *Client) Close() error {
	if !c.PoolByEnvID {
		return nil
	}

	var err error
	c.pool.Range(func(key, value interface{}) bool {
		err = value.(*memcache.Client).Close()
		return err == nil
	})
	return err
}
