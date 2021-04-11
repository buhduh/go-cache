package cache

import (
	"fmt"
	"time"
)

type cache struct {
	dataHandler DataHandler
	reaper      *Reaper
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
	dataHandler DataHandler,
	inv Invalidator,
	accessCB, createCB,
	updateCB ExtraCallback,
) Cacher {
	toRet := &cache{
		dataHandler: dataHandler,
		reaper:      NewReaper(inv, accessCB, createCB, updateCB),
		quit:        make(chan int8),
	}
	go toRet.begin()
	return toRet
}

type cacheElement struct {
	data     interface{}
	metadata Metadata
}

func (c *cache) Clear() {
	c.reaper.Clear()
	c.dataHandler.Clear()
}

func (c *cache) Put(key string, data interface{}) (interface{}, error) {
	found, err := c.dataHandler.Get(key)
	if err != nil {
		if !IsValueNotPresentError(err) {
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
			c.reaper.Update(&fCacheElem.metadata)
			toRet := fCacheElem.data
			fCacheElem.data = data
			c.dataHandler.Put(key, fCacheElem)
			return toRet, nil
		}
	}
	metadata := Metadata{}
	c.reaper.Create(&metadata)
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
		ok := IsValueNotPresentError(err)
		if ok {
			if len(data) == 1 {
				metadata := Metadata{}
				c.reaper.Create(&metadata)
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
	c.reaper.Access(&foundCacheElem.metadata)
	err = c.dataHandler.Put(key, foundCacheElem)
	return foundCacheElem.data, err
}

func (c *cache) Remove(key string) (interface{}, error) {
	found, err := c.dataHandler.Get(key)
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ValueNotPresentError{
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
	c.reaper.Remove()
	return cElem.data, nil
}

func (c *cache) Destroy() {
	c.dataHandler.Clear()
	c.reaper.Clear()
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
		if !c.reaper.IsValid(&elem.metadata) {
			c.dataHandler.Remove(key)
			c.reaper.Remove()
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
