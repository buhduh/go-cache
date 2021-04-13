package cache

import (
	"sync"
)

// NewInMemoryDataHandler returns a Datahandler that is backed with a sync.Map.
// This is the default DataHandler when nil is passed to NewCache.
func NewInMemoryDataHandler() DataHandler {
	return &inMemory{
		store: new(sync.Map),
	}
}

type inMemory struct {
	store *sync.Map
}

func (i *inMemory) Put(key string, data interface{}) error {
	i.store.Store(key, data)
	return nil
}

//Must throw a ValueNotPresentError and be nil
func (i *inMemory) Get(key string) (interface{}, error) {
	toRet, ok := i.store.Load(key)
	if !ok {
		return nil, ValueNotPresentError{
			Key: key,
		}
	}
	return toRet, nil
}

func (i *inMemory) Clear() error {
	//TODO, does this leak memory?
	i.store = new(sync.Map)
	return nil
}

//Must throw a ValueNotPresentError
//slightly less optimized, but the API is simpler, not sure if this is the
//right approach
func (i *inMemory) Remove(key string) error {
	_, ok := i.store.LoadAndDelete(key)
	if !ok {
		return ValueNotPresentError{
			Key: key,
		}
	}
	return nil
}

func (i *inMemory) Range(f func(string, interface{}) bool) {
	cb := func(iKey interface{}, val interface{}) bool {
		return f(iKey.(string), val)
	}
	i.store.Range(cb)
}
