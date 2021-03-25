package cache

import (
	"fmt"
	"sync"
	"time"
)

func NewInMemoryCache(lifetime *time.Duration) Cacher {
	quit := make(chan int8)
	cache := &inMemoryCache{
		store:    new(sync.Map),
		lifetime: lifetime,
		quit:     quit,
	}
	//sooooo weird, there's GOT to be a better way to do this...
	//circular reference...
	cache.reaper = newReaper(cache, quit)
	return cache
}

type Cacher interface {
	Put(string, interface{}) error
	Get(string, ...interface{}) (interface{}, error)
	Clear() error
	Remove(string) (interface{}, error)
}

type ValueNotPresentError struct {
	Key string
}

func (v ValueNotPresentError) Error() string {
	return fmt.Sprintf("no value found for key '%s'", v.Key)
}

type inMemoryCache struct {
	lifetime *time.Duration
	store    *sync.Map
	reaper   *reaperTracker
	quit     chan<- int8
}

type cacheElement struct {
	value  *interface{}
	ticker *time.Ticker
}

func (c *cacheElement) unPack() interface{} {
	return *(c.value)
}

func (i *inMemoryCache) newElem(data interface{}) *cacheElement {
	elem := new(cacheElement)
	if i.lifetime != nil {
		elem.ticker = time.NewTicker(*i.lifetime)
	}
	elem.value = &data
	return elem
}

func castCacheElement(element interface{}) (*cacheElement, bool) {
	toRet, ok := element.(*cacheElement)
	return toRet, ok
}

func (i *inMemoryCache) Put(key string, data interface{}) error {
	elem := i.newElem(data)
	i.store.Store(key, elem)
	go i.reaper.resetOrAddTimer(key, i.lifetime)
	return nil
}

func (i *inMemoryCache) Get(key string, data ...interface{}) (interface{}, error) {
	var rawElem interface{}
	if len(data) == 0 {
		rawElem, _ = i.store.Load(key)
	} else if len(data) == 1 {
		toStoreData := data[0]
		elem := i.newElem(toStoreData)
		rawElem, _ = i.store.LoadOrStore(key, elem)
	} else {
		return nil, fmt.Errorf(
			"only one item can be cached with InMemoryCache.Get(), attempted to store %d items",
			len(data),
		)
	}
	if rawElem == nil {
		return nil, ValueNotPresentError{Key: key}
	}
	var elem *cacheElement
	var ok bool
	if elem, ok = castCacheElement(rawElem); !ok {
		return nil, fmt.Errorf("could not cast cached value for key '%s' to a *cacheElement", key)
	}
	go i.reaper.resetOrAddTimer(key, i.lifetime)
	return elem.unPack(), nil
}

func (i *inMemoryCache) Clear() error {
	//I hope this actually recovers memory...
	i.store = new(sync.Map)
	close(i.quit)
	newQuit := make(chan int8)
	i.reaper = newReaper(i, newQuit)
	i.quit = newQuit
	return nil
}

func (i *inMemoryCache) Remove(key string) (interface{}, error) {
	rawElem, ok := i.store.LoadAndDelete(key)
	go i.reaper.remove(key)
	if !ok {
		return nil, ValueNotPresentError{Key: key}
	}
	var elem *cacheElement
	if elem, ok = castCacheElement(rawElem); !ok {
		return nil, fmt.Errorf("could not cast cached value for key '%s' to a *cachedElement", key)
	}
	return elem.unPack(), nil
}
