package cache

import (
	"time"
)

type Metadata struct {
	KeyCount *int64
	Created  int64
	Accessed int64
	Modified int64
	Extra    interface{}
}

type metadataHelper struct {
	keyCounter     chan int8
	quit           chan int8
	count          int64
	accessCallback ExtraCallback
	createCallback ExtraCallback
	updateCallback ExtraCallback
}

func newMetadataHelper(accessCB, createCB, updateCB ExtraCallback) *metadataHelper {
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
