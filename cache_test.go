package cache

import (
	"fmt"
	"github.com/buhduh/go-cache/datahandler"
	"github.com/buhduh/go-cache/invalidator"
	"sync"
	"testing"
)

type expCalls map[string]int

const (
	CREATE    string = "create"
	UPDATE           = "update"
	ACCESS           = "access"
	REMOVE           = "remove"
	CANCREATE        = "cancreate"
	ISVALID          = "isvalid"
	STOP             = "stop"
)

//breaking from the normal data loop for this, so much weird state i need to validate
func TestCache(t *testing.T) {
	inv := new(dummyInvalidator)
	handler := make(dummyHandler)
	tCache := NewCache(handler, inv)
	val, err := tCache.Put("foo", "yolo")
	if err != nil {
		t.Errorf("Put() should not have error'd, got '%s'", err)
	}
	if val != nil {
		t.Errorf("Put() should have returned nil for a new insert")
	}
	expCalls := map[string]int{
		CREATE:    0,
		UPDATE:    0,
		ACCESS:    0,
		REMOVE:    0,
		CANCREATE: 0,
		STOP:      0,
	}
	expCalls[CREATE] += 1
	expCalls[CANCREATE] += 1
	if msg := verifyCalls(expCalls, inv); msg != nil {
		t.Errorf("%s", *msg)
	}
	found, err := tCache.Get("foo")
	if err != nil || found == nil {
		if err != nil {
			t.Errorf("Get() should not have error'd, got '%s'", err)
		}
		if found == nil {
			t.Errorf("Get() should not have been nil")
		}
	} else {
		if res, ok := found.(string); !ok || res != "yolo" {
			if !ok {
				t.Errorf("could not convert cached element to a string")
			}
			if res != "yolo" {
				t.Errorf("expected '%s', got '%s'", "yolo", res)
			}
		}
	}
	expCalls[ACCESS] += 1
	if msg := verifyCalls(expCalls, inv); msg != nil {
		t.Errorf("%s", *msg)
	}
	found, err = tCache.Get("int", 42)
	if err != nil || found == nil {
		if err != nil {
			t.Errorf("Get() should not have error'd, got '%s'", err)
		}
		if found == nil {
			t.Errorf("Get() should not have returned nil")
		}
	}
	if res, ok := found.(int); !ok || res != 42 {
		if !ok {
			t.Errorf("Get() did not return an integer")
		} else if res != 42 {
			t.Errorf("Get() should have been 42")
		}
	}
	expCalls[CREATE] += 1
	expCalls[CANCREATE] += 1
	if msg := verifyCalls(expCalls, inv); msg != nil {
		t.Errorf("%s", *msg)
	}
	val, err = tCache.Put("int", 3)
	if err != nil {
		t.Errorf("Put() should not have error'd, got '%s'", err)
	}
	if tmp, ok := val.(int); !ok || tmp != 42 {
		if !ok {
			t.Errorf("Put() should have returned an integer")
		} else if tmp != 42 {
			t.Errorf("Put(), expected: %d, got %d", 42, tmp)
		}
	}
	expCalls[UPDATE] += 1
	if msg := verifyCalls(expCalls, inv); msg != nil {
		t.Errorf("%s", *msg)
	}
	removed, err := tCache.Remove("int")
	if err != nil {
		t.Errorf("failed removing an existing key, 'int'")
	}
	if val, ok := removed.(int); !ok || val != 3 {
		if !ok {
			t.Errorf("could not cast returned cache element to an interger")
		}
		if val != 3 {
			t.Errorf("expected 3, got '%d'", val)
		}
	}
	expCalls[REMOVE] += 1
	if msg := verifyCalls(expCalls, inv); msg != nil {
		t.Errorf("%s", *msg)
	}
	removed, err = tCache.Remove("int")
	if err == nil {
		t.Errorf("Remove() shouuld have returned an error, it did not")
	} else if !datahandler.IsValueNotPresentError(err) {
		t.Errorf("Remove() should have thrown a ValueNotPresentError, instead it threw, '%s'", err)
	}
	if removed != nil {
		t.Errorf("Remove() shoiuld have returend nil, it returned '%v'", removed)
	}
	for i := 0; i < 50; i++ {
		val, err := tCache.Put(fmt.Sprintf("key%d", i), i)
		if err != nil {
			t.Errorf("Put() should not have thrown an error, threw: '%s'", err)
		}
		if val != nil {
			t.Errorf("Put() should return nil for inserts, returned: '%v", val)
		}
		expCalls[CREATE] += 1
		expCalls[CANCREATE] += 1
		if msg := verifyCalls(expCalls, inv); msg != nil {
			t.Errorf("%s", *msg)
		}
	}
	val, err = tCache.Get("key20")
	if err != nil {
		t.Errorf("Put() should not have thrown an error, threw: '%s'", err)
	}
	if val != 20 {
		t.Errorf("Put(): expected %d, got %v", 20, val)
	}
	expCalls[ACCESS] += 1
	if msg := verifyCalls(expCalls, inv); msg != nil {
		t.Errorf("%s", *msg)
	}
	tCache.Clear()
	expCalls[STOP] += 1
	if len(handler) != 0 {
		t.Errorf("Clear() should have removed all elements, %d remaining", len(handler))
	}
}

type dummyInvalidator struct {
	numCreateCalled    int
	numUpdateCalled    int
	numAccessCalled    int
	numRemoveCalled    int
	numCanCreateCalled int
	numStopCalled      int
}

func (d *dummyInvalidator) Create(*invalidator.Metadata) {
	d.numCreateCalled += 1
}

func (d *dummyInvalidator) Update(*invalidator.Metadata) {
	d.numUpdateCalled += 1
}

func (d *dummyInvalidator) Access(*invalidator.Metadata) {
	d.numAccessCalled += 1
}

func (d *dummyInvalidator) Remove(*invalidator.Metadata) {
	d.numRemoveCalled += 1
}

func (d *dummyInvalidator) CanCreate(string, interface{}) bool {
	d.numCanCreateCalled += 1
	return true
}

func (d *dummyInvalidator) IsValid(*invalidator.Metadata) bool {
	return true
}

func (d *dummyInvalidator) Stop() {
	d.numStopCalled += 1
}

var mapLock = new(sync.Mutex)

type dummyHandler map[string]interface{}

func (d dummyHandler) Clear() error {
	for k, _ := range d {
		mapLock.Lock()
		delete(d, k)
		mapLock.Unlock()
	}
	return nil
}

func (d dummyHandler) Put(key string, data interface{}) error {
	mapLock.Lock()
	d[key] = data
	mapLock.Unlock()
	return nil
}

//Must throw a ValueNotPresentError and be nil
func (d dummyHandler) Get(key string) (interface{}, error) {
	toRet, ok := d[key]
	if !ok {
		return nil, datahandler.ValueNotPresentError{Key: key}
	}
	return toRet, nil
}

//Must throw a ValueNotPresentError and be bil
func (d dummyHandler) Remove(key string) error {
	if _, ok := d[key]; !ok {
		return datahandler.ValueNotPresentError{Key: key}
	}
	mapLock.Lock()
	delete(d, key)
	mapLock.Unlock()
	return nil
}

func (d dummyHandler) Range(f func(string, interface{}) bool) {
	mapLock.Lock()
	defer mapLock.Unlock()
	for k, v := range d {
		if !f(k, v) {
			return
		}
	}
}

func verifyCalls(calls map[string]int, inv *dummyInvalidator) *string {
	msg := func(f string, exp, res int) *string {
		toRet := fmt.Sprintf("func: '%s', expected %d, got %d", f, exp, res)
		return &toRet
	}
	if calls[CREATE] != inv.numCreateCalled {
		return msg("create", calls[CREATE], inv.numCreateCalled)
	}
	if calls[UPDATE] != inv.numUpdateCalled {
		return msg("update", calls[UPDATE], inv.numUpdateCalled)
	}
	if calls[ACCESS] != inv.numAccessCalled {
		return msg("access", calls[ACCESS], inv.numAccessCalled)
	}
	if calls[REMOVE] != inv.numRemoveCalled {
		return msg("remove", calls[REMOVE], inv.numRemoveCalled)
	}
	if calls[CANCREATE] != inv.numCanCreateCalled {
		return msg("canCreate", calls[CANCREATE], inv.numCanCreateCalled)
	}
	if calls[STOP] != inv.numStopCalled {
		return msg("stop", calls[STOP], inv.numStopCalled)
	}
	return nil
}
