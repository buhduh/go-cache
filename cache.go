package cache

import (
	"fmt"
	"sync"
)

type Cacher interface {
	Put(string, interface{}) error
	Get(string, ...interface{}) (interface{}, error)
	Clear() error
	Remove(string) (interface{}, error)
}

type StringCacher interface {
	Put(string, string) error
	Get(string, ...string) (string, error)
	Clear() error
	Remove(string) (string, error)
}

type inMemoryCache struct {
	store *sync.Map
}

type inMemoryStringCache inMemoryCache

type ValueNotPresentError struct {
	Key string
}

func (v ValueNotPresentError) Error() string {
	return fmt.Sprintf("no value found for key '%s'", v.Key)
}

func NewInMemoryCacher() Cacher {
	return &inMemoryCache{
		store: new(sync.Map),
	}
}

func NewInMemoryStringCacher() StringCacher {
	return &inMemoryStringCache{
		store: new(sync.Map),
	}
}

func (i *inMemoryCache) Put(key string, data interface{}) error {
	i.store.Store(key, data)
	return nil
}

func (i *inMemoryCache) Get(key string, data ...interface{}) (interface{}, error) {
	var toRet interface{}
	if len(data) == 0 {
		toRet, _ = i.store.Load(key)
	} else if len(data) == 1 {
		toRet, _ = i.store.LoadOrStore(key, data[0])
	} else {
		return nil, fmt.Errorf(
			"only one item can be cached with InMemoryCache.Get(), attempted to store %d items",
			len(data),
		)
	}
	if toRet == nil {
		return nil, ValueNotPresentError{Key: key}
	}
	return toRet, nil
}

func (i *inMemoryCache) Clear() error {
	//I hope this actually recovers memory...
	i.store = new(sync.Map)
	return nil
}

func (i *inMemoryCache) Remove(key string) (interface{}, error) {
	toRet, ok := i.store.LoadAndDelete(key)
	if !ok {
		return nil, ValueNotPresentError{Key: key}
	}
	return toRet, nil
}

func (i *inMemoryStringCache) Put(key string, value string) error {
	return ((*inMemoryCache)(i)).Put(key, value)
}

func (i *inMemoryStringCache) Get(key string, value ...string) (string, error) {
	vals := make([]interface{}, len(value))
	for i, v := range value {
		vals[i] = v
	}
	toRet, err := ((*inMemoryCache)(i)).Get(key, vals...)
	if err != nil {
		return "", err
	}
	return toRet.(string), nil
}

func (i *inMemoryStringCache) Clear() error {
	//I hope this actually recovers memory...
	i.store = new(sync.Map)
	return nil
}

func (i *inMemoryStringCache) Remove(key string) (string, error) {
	toRet, err := ((*inMemoryCache)(i)).Remove(key)
	if err != nil {
		return "", err
	}
	return toRet.(string), nil
}
