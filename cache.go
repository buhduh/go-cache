// Package cache provides a modifiable thread safe cache system that separates concerns along the way data
// is stored in the backend with the DataHandler interface and when/how data is validated
// with the Invalidator interface.
package cache

import (
	"fmt"
	"time"
)

// Cacher primary interface for this package.
type Cacher interface {
	// Clear remove all elements from the cache.
	Clear()
	// Get a single element from the cache, if a second parameter is
	// provided, will set the cache to that value if nothing is present
	// returns a ValueNotPresentError if no value was found at key.
	Get(string, ...interface{}) (interface{}, error)
	// Put a value at key
	Put(string, interface{}) (interface{}, error)
	// Remove a single item, returning the item or a ValueNotPresentError
	// if no item is present.
	Remove(string) (interface{}, error)
	// Destroy the cache releasing resources.
	Destroy()
}

// NewCache returns a Cacher Interface whose behavior is determined by
// datahandler and inv.
func NewCache(
	// if nil, an InMemoryDataHandler will be used.
	dataHandler DataHandler,
	// if nil, a nopInvalidator will be used that maintains
	// Metadata consistently and IsValid always returns true.
	inv Invalidator,
) Cacher {
	if inv == nil {
		inv = &nopInvalidator{}
	}
	if dataHandler == nil {
		dataHandler = NewInMemoryDataHandler()
	}
	toRet := &cache{
		dataHandler: dataHandler,
		reaper:      newReaper(inv),
		quit:        make(chan int8),
	}
	go toRet.begin()
	return toRet
}

// DataHandler is the interface that the Cacher interface uses to
// actually store data to the cache.
type DataHandler interface {
	// Put a single item in the cache
	Put(string, interface{}) error
	// Get a sinle item from the cache, must return a ValueNotPresentError
	// if there is nothing in the cache there.
	Get(string) (interface{}, error)
	// Clear removes all elements from the cache.
	Clear() error
	// Remove an item from the cache, returns a ValueNotPresentError
	// if there was nothing in the cache at key.
	Remove(string) error
	// Range iterates through all items in the cache and calls
	// the passed in function.  If the function returns false, iteration halts.
	Range(func(string, interface{}) bool)
}

// Invalidator is the interface that Cacher uses to determine if an item is valid.
// If IsValid returns false, the item will be removed from the cache.
type Invalidator interface {
	// IsValid determines whether or not a cache item is valid.
	// This package constantly loops through the cache in a background go routine
	// calling IsValid on all items.
	IsValid(*Metadata) bool

	// The following functions provide a means for implentations of this interface
	// to modify Metadata in ways not predicted by this package. Intended to modify
	// the 'Extra' field in Metadata.

	// AccessExtra is called whenever Cacher.Get() is called and an item
	// is present.
	AccessExtra(*Metadata)
	// CreateExtra is called when Cacher.Put or in the case of an insertion when Cacher.Get
	// is called.
	CreateExtra(*Metadata)
	// UpdateExtra is called whenever Cacher.Put overwrites an existing item.
	UpdateExtra(*Metadata)
}

// ValueNotPresentError is returned when an item isn't found in the cache.
type ValueNotPresentError struct {
	Key string // The item key.
}

// Error satisfies the Error interface.
func (v ValueNotPresentError) Error() string {
	return fmt.Sprintf("no value found for key '%s'", v.Key)
}

// IsValueNotPresentError is a simple test to determine if an error
// is of type 'ValueNotPresentError'.
func IsValueNotPresentError(err error) bool {
	_, ok := err.(ValueNotPresentError)
	return ok
}

type nopInvalidator struct{}

func (n *nopInvalidator) IsValid(*Metadata) bool {
	return true
}

func (n *nopInvalidator) AccessExtra(*Metadata) {}
func (n *nopInvalidator) CreateExtra(*Metadata) {}
func (n *nopInvalidator) UpdateExtra(*Metadata) {}

type reaper struct {
	*metadataHelper
	Invalidator
}

func newReaper(inv Invalidator) *reaper {
	return &reaper{
		newMetadataHelper(inv.AccessExtra, inv.CreateExtra, inv.UpdateExtra),
		inv,
	}
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
	dur, _ := time.ParseDuration("100ms")
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

type cache struct {
	dataHandler DataHandler
	reaper      *reaper
	quit        chan int8
}
