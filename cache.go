package cache

import (
	"fmt"
	"github.com/buhduh/go-cache/datahandler"
	"github.com/buhduh/go-cache/invalidator"
	"time"
)

type cache struct {
	dataHandler datahandler.DataHandler
	invalidator invalidator.Invalidator
	quit        chan int8
}

type Cacher interface {
	Clear()
	Get(string, ...interface{}) (interface{}, error)
	Put(string, interface{}) (interface{}, error)
	Remove(string) (interface{}, error)
	Destroy()
}

func NewCache(
	dataHandler datahandler.DataHandler,
	invalidator invalidator.Invalidator,
) Cacher {
	toRet := &cache{
		dataHandler: dataHandler,
		invalidator: invalidator,
		quit:        make(chan int8),
	}
	go toRet.begin()
	return toRet
}

type cacheElement struct {
	data     interface{}
	metadata invalidator.Metadata
}

func (c *cache) Clear() {
	c.invalidator.Stop()
	c.dataHandler.Clear()
}

func (c *cache) Put(key string, data interface{}) (interface{}, error) {
	found, err := c.dataHandler.Get(key)
	if err != nil {
		if !datahandler.IsValueNotPresentError(err) {
			return nil, err
		}
	} else if found != nil {
		fCacheElem, ok := found.(cacheElement)
		if !ok {
			return nil, fmt.Errorf(
				"cache may be corrupt, found something for key '%s', but can't unpack it",
				key,
			)
		} else {
			c.invalidator.Update(&fCacheElem.metadata)
			toRet := fCacheElem.data
			fCacheElem.data = data
			c.dataHandler.Put(key, fCacheElem)
			return toRet, nil
		}
	}
	if can := c.invalidator.CanCreate(key, data); !can {
		return nil, fmt.Errorf(
			"unabled to create element with key: '%s', and value: '%v'",
			key, data,
		)
	}
	metadata := invalidator.Metadata{}
	c.invalidator.Create(&metadata)
	err = c.dataHandler.Put(
		key,
		cacheElement{
			data:     data,
			metadata: metadata,
		},
	)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *cache) Get(key string, data ...interface{}) (interface{}, error) {
	if len(data) > 1 {
		return nil, fmt.Errorf(
			"only a single value can be sent to Get to be cached as a default, attemped to pass %d items",
			len(data),
		)
	}
	found, err := c.dataHandler.Get(key)
	if err != nil {
		//new element
		ok := datahandler.IsValueNotPresentError(err)
		if ok {
			if len(data) == 1 {
				if can := c.invalidator.CanCreate(key, data[0]); !can {
					return nil, fmt.Errorf(
						"unable to create element with key, '%s' and value '%v'",
						key, data[0],
					)
				}
				metadata := invalidator.Metadata{}
				c.invalidator.Create(&metadata)
				putErr := c.dataHandler.Put(
					key,
					cacheElement{
						data:     data[0],
						metadata: metadata,
					},
				)
				if putErr != nil {
					return nil, putErr
				} else {
					return data[0], nil
				}
			} else {
				return nil, err
			}
		}
	}
	foundCacheElem, ok := found.(cacheElement)
	if !ok {
		return nil, fmt.Errorf(
			"cache may be corrupt, found something for key '%s', but can't unpack it",
			key,
		)
	}
	c.invalidator.Access(&foundCacheElem.metadata)
	err = c.dataHandler.Put(key, foundCacheElem)
	return foundCacheElem.data, err
}

func (c *cache) Remove(key string) (interface{}, error) {
	found, err := c.dataHandler.Get(key)
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, datahandler.ValueNotPresentError{
			Key: key,
		}
	}
	err = c.dataHandler.Remove(key)
	if err != nil {
		return nil, err
	}
	cElem, ok := found.(cacheElement)
	if !ok {
		return nil, fmt.Errorf(
			"cache may be corrupt, found something for key '%s', but can't unpack it",
			key,
		)
	}
	c.invalidator.Remove(&cElem.metadata)
	return cElem.data, nil
}

func (c *cache) Destroy() {
	c.dataHandler.Clear()
	c.invalidator.Stop()
	close(c.quit)
}

func (c *cache) begin() {
	cb := func(key string, val interface{}) bool {
		select {
		case <-c.quit:
			return false
		default:
		}
		elem, ok := val.(cacheElement)
		if !ok {
			return true
		}
		if !c.invalidator.IsValid(&elem.metadata) {
			c.dataHandler.Remove(key)
		}
		return true
	}
	dur, _ := time.ParseDuration("500ms")
	myTicker := time.NewTicker(dur)
	for {
		select {
		case <-c.quit:
			myTicker.Stop()
			return
		case <-myTicker.C:
			c.dataHandler.Range(cb)
		}
	}
}
