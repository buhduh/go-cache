package cache

import (
	"sync"
	"time"
)

type reaperTracker struct {
	timers map[string]*time.Timer
	quit   <-chan int8
	lock   *sync.Mutex
}

func newReaper(cache Cacher, quit <-chan int8) *reaperTracker {
	reaper := &reaperTracker{
		timers: make(map[string]*time.Timer),
		quit:   quit,
		lock:   new(sync.Mutex),
	}
	go reaper.run(cache)
	return reaper
}

func (r *reaperTracker) run(cache Cacher) {
	for {
		for k, t := range r.timers {
			select {
			case <-t.C:
				cache.Remove(k)
				//according to the docs this is safe...
				r.lock.Lock()
				delete(r.timers, k)
			case <-r.quit:
				for _, toStop := range r.timers {
					toStop.Stop()
				}
			default:
			}
		}
	}
}

func (r *reaperTracker) resetOrAddTimer(key string, lifetime *time.Duration) {
	if lifetime == nil {
		return
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	if timer, ok := r.timers[key]; ok {
		timer.Reset(*lifetime)
	} else {
		r.timers[key] = time.NewTimer(*lifetime)
	}
}

func (r *reaperTracker) remove(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if t, ok := r.timers[key]; ok {
		t.Stop()
		delete(r.timers, key)
	}
}
