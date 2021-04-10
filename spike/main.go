package main

import (
	"fmt"
	cache "github.com/buhduh/go-cache"
	"github.com/buhduh/go-cache/datahandler"
	"github.com/buhduh/go-cache/invalidator"
	"time"
)

func main() {
	exp, _ := time.ParseDuration("1s")
	myCache := cache.NewCache(
		datahandler.NewInMemoryDataHandler(),
		invalidator.NewTimedInvalidator(exp),
	)
	myCache.Put("key", "key")
	key, _ := myCache.Put("key", "foo")
	fmt.Printf("should be 'key': '%s'\n", key.(string))
	key, _ = myCache.Get("key")
	fmt.Printf("should be 'foo': '%s'\n", key.(string))
	twoSec, _ := time.ParseDuration("2s")
	time.Sleep(twoSec)
	gone, err := myCache.Get("key")
	fmt.Printf("should be not found: '%s'", err)
	fmt.Printf("should be nil: %v", gone)
}
