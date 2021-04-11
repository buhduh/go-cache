package cache

import (
	"time"
)

// NewTimedInvalidator returns an Invaidator that validates cache based on lifetime.
// Takes the most recent value of Metadata.Accessed, Metadata.Created, or Metadata.Updated
// and compares to lifefime.
func NewTimedInvalidator(lifetime time.Duration) Invalidator {
	return &timedInvalidator{
		lifetime: lifetime,
	}
}

// CanCreate always returns true.
func (t *timedInvalidator) CanCreate(*Metadata, string, interface{}) bool {
	return true
}

// IsValid compares the most recent of Metadata.Accessed, Metadata.Created, or Metadata.Updated
// and lifetime.
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

func (t *timedInvalidator) AccessExtra(*Metadata) {}
func (t *timedInvalidator) CreateExtra(*Metadata) {}
func (t *timedInvalidator) UpdateExtra(*Metadata) {}

type timedInvalidator struct {
	lifetime time.Duration
}
