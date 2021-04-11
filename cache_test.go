package cache

import (
//"testing"
)

/*
func TestCache(t *testing.T) {
	//huh?
	dur, _ := time.ParseDuration("1s")
	myCache := NewCache(
		NewInMemoryDataHandler(),
		NewTimedInvalidator(dur),
		nil, nil, nil,
	)
	got, err := myCache.Get("foo", 42)
	if err != nil {
		t.Errorf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		t.Errorf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 42 {
		t.Errorf("shold have gotten 42")
	}
	got, err = myCache.Get("foo", 3)
	if err != nil {
		t.Errorf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		t.Errorf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 42 {
		t.Errorf("shold have gotten 42")
	}
	got, err = myCache.Remove("foo")
	if err != nil {
		t.Errorf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		t.Errorf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 42 {
		t.Errorf("shold have gotten 42")
	}
	got, err = myCache.Put("foo", 3)
	if err != nil {
		t.Errorf("shouldn't have error'd, got '%s'\n", err)
	}
	if got != nil {
		t.Errorf("should be nil\n")
	}
	got, err = myCache.Put("foo", 10)
	if err != nil {
		t.Errorf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		t.Errorf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 3 {
		t.Errorf("shold have gotten 3")
	}
	got, err = myCache.Get("foo", 3)
	if err != nil {
		t.Errorf("shouldn't have error'd, got '%s'\n", err)
	}
	if got == nil {
		t.Errorf("result shouldn't be nil")
	} else if res, ok := got.(int); !ok || res != 10 {
		t.Errorf("shold have gotten 10")
	}
	oneSec, _ := time.ParseDuration("2s")
	time.Sleep(oneSec)
	got, err = myCache.Get("foo")
	if err != nil && !IsValueNotPresentError(err) {
		t.Errorf("Get should return ValueNotPresentError")
	}
	if got != nil {
		t.Errorf("Get() should be nil: %v", got)
	}
}
*/

func ExampleNewInMemoryDataHandler() {
	// InMemoryDataHandler is the default.
	// could also do: NewCache(NewInMemoryDataHandler(), nil)
	myCache := NewCache(nil, nil)
	newItem, _ := myCache.Get("foo", 42)
	found, _ := myCache.Get("foo")
	if newItem.(int) != found.(int) {
		println("this won't happen")
	}
}
