package cache

import (
	"time"
)

// Metadata is primarily used by Invalidator to determine if
// a cache item is valid.  Invalidator.AccessExtra, Invalidator.CreateExtra, and Invalidator.UpdateExtra
// are intended to modify the Extra field in Metadata.
type Metadata struct {
	// KeyCount is a thread safe pointer to the total count of the cache
	// KeyCount is managed in a background go routine.
	KeyCount *int64
	// Created is a Unix time stamp when an item was originally inserted into the cache.
	Created int64
	// Accessed is a Unix time stamp of the last time an item was retrieved with Cacher.Get
	Accessed int64
	// Modified is a Unix time stamp of the last time an item was modfied with Cacher.Put
	Modified int64
	// Extra provides a means for an outside implementation of Invalidator to determine
	// if an item is valid.
	Extra interface{}
}

type metadataHelper struct {
	keyCounter     chan int8
	quit           chan int8
	count          int64
	accessCallback func(*Metadata)
	createCallback func(*Metadata)
	updateCallback func(*Metadata)
}

func newMetadataHelper(accessCB, createCB, updateCB func(*Metadata)) *metadataHelper {
	toRet := &metadataHelper{
		keyCounter:     make(chan int8),
		quit:           make(chan int8),
		count:          0,
		accessCallback: accessCB,
		createCallback: createCB,
		updateCallback: updateCB,
	}
	go toRet.begin()
	return toRet
}

func (m *metadataHelper) Create(data *Metadata) {
	m.keyCounter <- 1
	data.Created = time.Now().Unix()
	data.KeyCount = m.getCount()
	if m.createCallback != nil {
		m.createCallback(data)
	}
}

func (m *metadataHelper) Access(data *Metadata) {
	data.Accessed = time.Now().Unix()
	if m.accessCallback != nil {
		m.accessCallback(data)
	}
}

func (m *metadataHelper) Update(data *Metadata) {
	data.Modified = time.Now().Unix()
	if m.updateCallback != nil {
		m.updateCallback(data)
	}
}

func (m *metadataHelper) Remove() {
	m.keyCounter <- -1
}

func (m *metadataHelper) Clear() {
	close(m.quit)
	m.count = 0
	m.quit = make(chan int8)
	go m.begin()
}

func (m *metadataHelper) getCount() *int64 {
	return &m.count
}

func (m *metadataHelper) begin() {
	for {
		select {
		case tmp := <-m.keyCounter:
			m.count += int64(tmp)
		case <-m.quit:
			return
		}
	}
}
