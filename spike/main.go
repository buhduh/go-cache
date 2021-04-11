package main

import (
	cache "github.com/buhduh/go-cache"
)

func main() {
	reaper := cache.NewReaper(cache.NopInvalidator{})
	reaper := new(cache.Metadata)
	reaper.Create(mData)
}
