package main

import (
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	uri := "localhost:11211"
	mc := memcache.New(uri)

	err := mc.Set(&memcache.Item{
		Key:   "picop-example",
		Value: []byte("example"),
	})
	if err != nil {
		log.Fatal(err)
	}

	it, err := mc.Get("picop-example")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(it.Value))
}
