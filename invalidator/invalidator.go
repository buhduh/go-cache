package invalidator

import (
	"time"
)

type Invalidator interface {
	Create(*Metadata)
	Update(*Metadata)
	Access(*Metadata)
	Remove(*Metadata)
	CanCreate(string, interface{}) bool
	IsValid(*Metadata) bool
	Stop()
}

type Metadata struct {
	KeyCount *int64
	Created  int64
	Accessed int64
	Modified int64
	Extra    interface{}
}

func createHelper(m *Metadata) {
	m.Created = time.Now().Unix()
	m.Accessed = -1
	m.Modified = -1
}

func updateHelper(m *Metadata) {
	m.Modified = time.Now().Unix()
}

func accessHelper(m *Metadata) {
	m.Accessed = time.Now().Unix()
}
