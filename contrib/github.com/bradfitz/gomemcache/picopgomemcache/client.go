package picopgomemcache

import (
	otelprop "go.opentelemetry.io/otel/propagation"

	// Please use github.com/picop-rd/gomemcache instead of github.com/bradfitz/gomemcache/memcache by replacing in go.mod
	"github.com/bradfitz/gomemcache/memcache"
)

type Client struct {
	server     []string
	propagator otelprop.TextMapPropagator
}

func New(propagator otelprop.TextMapPropagator, server ...string) *Client {
	return &Client{
		server:     server,
		propagator: propagator,
	}
}

func (c *Client) Connect() *memcache.Client {
	mc := memcache.New(c.server...)
	mc.DialContext = DialContext(c.propagator)
	return mc
}
