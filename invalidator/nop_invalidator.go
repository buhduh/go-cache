package invalidator

import (
	"github.com/buhduh/go-cache/counter"
)

type NopInvalidator struct {
	counter *counter.Counter
}

func NewNopInvalidator() *NopInvalidator {
	return &NopInvalidator{
		counter: counter.NewCounter(),
	}
}

func (n *NopInvalidator) Create(data *Metadata) {
	n.counter.Update(1)
	data.KeyCount = n.counter.Get()
}

func (n *NopInvalidator) Update(*Metadata) {
}

func (n *NopInvalidator) Access(*Metadata) {
}

func (n *NopInvalidator) Remove(*Metadata) {
	n.counter.Update(-1)
}

func (n *NopInvalidator) CanCreate(string, *Metadata) bool {
	return true
}

func (n *NopInvalidator) IsValid(*Metadata) bool {
	return true
}

func (n *NopInvalidator) Stop() {
}
