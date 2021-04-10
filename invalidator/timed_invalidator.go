package invalidator

import (
	"github.com/buhduh/go-cache/counter"
	"time"
)

type timedInvalidator struct {
	lifetime time.Duration
	counter  *counter.Counter
}

func NewTimedInvalidator(lifetime time.Duration) Invalidator {
	return &timedInvalidator{
		lifetime: lifetime,
		counter:  counter.NewCounter(),
	}
}

func (t *timedInvalidator) Create(data *Metadata) {
	createHelper(data)
	t.counter.Update(1)
	data.KeyCount = t.counter.Get()
}

func (t *timedInvalidator) Update(data *Metadata) {
	updateHelper(data)
}

func (t *timedInvalidator) Access(data *Metadata) {
	accessHelper(data)
}

func (t *timedInvalidator) Remove(data *Metadata) {
	t.counter.Update(-1)
}

func (t *timedInvalidator) CanCreate(string, interface{}) bool {
	return true
}

func (t *timedInvalidator) IsValid(data *Metadata) bool {
	var max int64
	if data.Created > data.Accessed {
		max = data.Created
	} else {
		max = data.Accessed
	}
	if data.Modified > max {
		max = data.Modified
	}
	return max >= time.Now().Add(-1*t.lifetime).Unix()
}

//do i need this?
func (t *timedInvalidator) Stop() {
}
