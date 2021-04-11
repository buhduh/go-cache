package main

import cache "github.com/buhduh/go-cache"

func main() {
	myCache := cache.NewCache(nil, nil)
	newItem, _ := myCache.Get("foo", 42)
	found, _ := myCache.Get("foo")
	if newItem.(int) != found.(int) {
		println("not equal")
	} else {
		println("equal")
	}
}
