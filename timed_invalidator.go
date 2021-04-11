package cache

import (
	"time"
)

type timedInvalidator struct {
	lifetime time.Duration
}

func NewTimedInvalidator(lifetime time.Duration) Invalidator {
	return &timedInvalidator{
		lifetime: lifetime,
	}
}

func (t *timedInvalidator) CanCreate(*Metadata, string, interface{}) bool {
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
