package main

import (
	"fmt"
	cache "github.com/buhduh/go-cache"
	"time"
)

func main() {
	dur, _ := time.ParseDuration("1s")
	myCache := cache.NewCache(
		cache.NewInMemoryDataHandler(),
		cache.NewTimedInvalidator(dur),
		nil, nil, nil,
	)
	got, err := myCache.Get("foo", 42)
	if err != nil {
		fmt.Printf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		fmt.Printf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 42 {
		fmt.Printf("shold have gotten 42")
	}
	fmt.Printf("got: %v\n", got)
	got, err = myCache.Get("foo", 3)
	if err != nil {
		fmt.Printf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		fmt.Printf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 42 {
		fmt.Printf("shold have gotten 42")
	}
	fmt.Printf("got: %v\n", got)
	got, err = myCache.Remove("foo")
	if err != nil {
		fmt.Printf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		fmt.Printf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 42 {
		fmt.Printf("shold have gotten 42")
	}
	got, err = myCache.Put("foo", 3)
	if err != nil {
		fmt.Printf("shouldn't have error'd, got '%s'\n", err)
	}
	if got != nil {
		fmt.Printf("should be nil\n")
	}
	got, err = myCache.Put("foo", 10)
	if err != nil {
		fmt.Printf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		fmt.Printf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 3 {
		fmt.Printf("shold have gotten 3")
	}
	got, err = myCache.Get("foo", 3)
	if err != nil {
		fmt.Printf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		fmt.Printf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 10 {
		fmt.Printf("shold have gotten 10")
	}
	oneSec, _ := time.ParseDuration("2s")
	time.Sleep(oneSec)
	got, err = myCache.Get("foo")
	if err != nil && !cache.IsValueNotPresentError(err) {
		println("1")
		fmt.Printf("Get should return ValueNotPresentError")
	}
	if got != nil {
		fmt.Printf("Get() should be nil")
	}
}
