package cache

import (
	"fmt"
	"testing"
	"time"
)

// forgoing the usual testCases = []*struct{}{} style here
// could not come up with an adequate solution in this format,
// just going to write a bunch of LONG cache calls checking for
// expected reuslts along the way.
func TestCache(t *testing.T) {
	t.Run("method=Get", testGet)
	t.Run("method=Put", testPut)
}

type simple struct {
	key string
	bar string
}

func isFooEqual(f1, f2 simple) bool {
	return (f1.key == f2.key) && (f1.bar == f2.bar)
}

type complexVal struct {
	pString *string
	mString string
	pInt    *int
	mInt    int
	pFoo    *simple
	mFoo    simple
}

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}

func isComplexValEqual(c1, c2 complexVal) bool {
	if !(c1.pString == nil) == (c2.pString == nil) {
		return false
	}
	if c1.pString != nil && *c1.pString != *c2.pString {
		return false
	}
	if c1.mString != c2.mString {
		return false
	}
	if !(c1.pInt == nil) == (c2.pInt == nil) {
		return false
	}
	if c1.pInt != nil && *c1.pInt != *c2.pInt {
		return false
	}
	if c2.mInt != c2.mInt {
		return false
	}
	if !(c1.pFoo == nil) == (c2.pFoo == nil) {
		return false
	}
	if c1.pFoo != nil && !isFooEqual(*c1.pFoo, *c2.pFoo) {
		return false
	}
	return isFooEqual(c1.mFoo, c2.mFoo)
}

type dummyInvalidator struct {
	accessCalled int
	createCalled int
	updateCalled int
	lastMetadata Metadata
}

func (d *dummyInvalidator) IsValid(*Metadata) bool {
	return true
}

func (d *dummyInvalidator) AccessExtra(data *Metadata) {
	d.accessCalled += 1
	d.lastMetadata = *data
}

func (d *dummyInvalidator) CreateExtra(data *Metadata) {
	d.createCalled += 1
	d.lastMetadata = *data
}

func (d *dummyInvalidator) UpdateExtra(data *Metadata) {
	d.updateCalled += 1
	d.lastMetadata = *data
}

func (d *dummyInvalidator) checkExtras(state []int) bool {
	return d.accessCalled == state[0] &&
		d.createCalled == state[1] &&
		d.updateCalled == state[2]
}

func (d *dummyInvalidator) getMetadata() Metadata {
	return d.lastMetadata
}

func (d *dummyInvalidator) checkMetadata(data Metadata) bool {
	return d.lastMetadata.Accessed == data.Accessed &&
		d.lastMetadata.Created == data.Created &&
		d.lastMetadata.Modified == data.Modified
}

func testPut(t *testing.T) {
	handler := make(dummyHandler)
	invalidator := new(dummyInvalidator)
	myCache := NewCache(handler, invalidator)
	cpxVal1 := complexVal{
		pString: String("foo"),
		mString: "bar",
		pInt:    Int(42),
		mInt:    32,
		pFoo: &simple{
			key: "key",
			bar: "bar",
		},
		mFoo: simple{
			key: "huh",
			bar: "what",
		},
	}
	val, err := myCache.Put("foo", cpxVal1)
	expData := Metadata{
		Accessed: -1,
		Created:  time.Now().Unix(),
		Modified: -1,
	}
	if !invalidator.checkExtras([]int{0, 1, 0}) {
		t.Errorf("Cacher.Put() extra function calls inconsistent")
	}
	if !invalidator.checkMetadata(expData) {
		t.Errorf("Cacher.Put() metadata is incorrect")
	}
	if err != nil {
		t.Errorf(
			"Cacher.Put() should not have returned an error, returned error: '%s'", err,
		)
	}
	if val != nil {
		t.Errorf(
			"Cacher.Put() initial insert of item should return nil, returned: '%#v\n",
			val,
		)
	}
	var found interface{}
	var ok bool
	found, ok = handler["foo"]
	if !ok {
		t.Errorf("Cacher.Put() did not insert item into datahandler")
	}
	cElem, ok := found.(cacheElement)
	if !ok {
		t.Errorf("Cacher.Put() did not insert a cacheElement")
	}
	resCPX, ok := (cElem.data).(complexVal)
	if !ok {
		t.Errorf("Cacher.Put() could not unpack cached value")
	}
	if !isComplexValEqual(cpxVal1, resCPX) {
		t.Errorf(
			"Cacher.Put() modified something when writing, expected '%#v, got '%#v'",
			cpxVal1, resCPX,
		)
	}
	cpxVal2 := &complexVal{
		pString: String("a pointer to a string"),
		mString: "a string member",
		pInt:    Int(43),
		mInt:    33,
		pFoo: &simple{
			key: "this is a string",
			bar: "this is another string",
		},
		mFoo: simple{
			key: "yolo",
			bar: "asdfadsf",
		},
	}
	expData.Modified = time.Now().Unix()
	resItem, err := myCache.Put("foo", cpxVal2)
	if !invalidator.checkExtras([]int{0, 1, 1}) {
		t.Errorf("Cacher.Put() extra function calls inconsistent")
	}
	if !invalidator.checkMetadata(expData) {
		t.Errorf("Cacher.Put() metadata is incorrect")
	}
	if err != nil {
		t.Errorf("Cacher.Put() should not have error'd, got '%s'", err)
	}
	if resItem == nil {
		t.Fatalf("Cacher.Put() returned a nil item, cannot continue")
	}
	if resCPX, ok = resItem.(complexVal); !ok {
		t.Errorf("Cacher.Put() did not return a complexVal, returned: '%#v'", resCPX)
	}
	if !isComplexValEqual(cpxVal1, resCPX) {
		t.Errorf(
			"Cacher.Put() did not return the previously cached item, expected: '%#v\n', got: '%#v'\n",
			cpxVal1, resCPX,
		)
	}
	cItem, ok := handler["foo"]
	if !ok {
		t.Errorf("Cacher.Put() did not update the cache at 'foo'")
	} else {
		cElem, ok := cItem.(cacheElement)
		if !ok {
			t.Errorf("Cacher.Put() did not save a cachedElement")
		} else {
			var cCPXVal *complexVal
			if cCPXVal, ok = (cElem.data).(*complexVal); !ok || cCPXVal == nil {
				t.Errorf("Cacher.Put() did not save a pointer to a complexVal")
			} else if !isComplexValEqual(*cpxVal2, *cCPXVal) {
				t.Errorf(
					"Cacher.Put() did not cache the correct item, expected '%#v', got: '%#v'\n",
					cpxVal2, *cCPXVal,
				)
			}
		}
	}
}

func testGet(t *testing.T) {
	handler := make(dummyHandler)
	invalidator := new(dummyInvalidator)
	myCache := NewCache(handler, invalidator)
	expData := Metadata{
		Accessed: -1,
		Created:  time.Now().Unix(),
		Modified: -1,
	}
	foo, err := myCache.Get("foo", "bar")
	if !invalidator.checkExtras([]int{0, 1, 0}) {
		t.Errorf("Cacher.Put() extra function calls inconsistent")
	}
	if !invalidator.checkMetadata(expData) {
		t.Errorf("Cacher.Put() metadata is incorrect")
	}
	if err != nil {
		t.Errorf("Caaher.Get() should not have error'd, got '%s'", err)
	}
	if val, ok := foo.(string); !ok || val != "bar" {
		if !ok {
			t.Errorf("Cacher.Get() did not return the item when inserting")
		} else if val != "bar" {
			t.Errorf(
				"Cacher.Get() did not return correct value, expected '%s', got '%s",
				"bar", val,
			)
		}
	}
	handler["bar"] = cacheElement{
		data:     "baz",
		metadata: Metadata{},
	}
	expData.Accessed = time.Now().Unix()
	bar, err := myCache.Get("bar")
	if !invalidator.checkExtras([]int{1, 1, 0}) {
		t.Errorf("Cacher.Put() extra function calls inconsistent")
	}
	if !invalidator.checkMetadata(expData) {
		t.Errorf(
			"Cacher.Put() metadata is incorrect, expected:\n%s\ngot:\n%s\n",
			expData, invalidator.getMetadata(),
		)
	}
	if err != nil {
		t.Errorf("Cacher.Get() should not have error'd, got '%s'", err)
	}
	if val, ok := bar.(string); !ok || val != "baz" {
		if !ok {
			t.Errorf("Cacher.Get() did not return the expected type")
		} else if val != "baz" {
			t.Errorf(
				"Cacher.Get() did nor return expected value, expected '%s', got '%s'",
				"baz", val,
			)
		}
	}
	val, err := myCache.Get("yolo")
	if !invalidator.checkExtras([]int{1, 1, 0}) {
		t.Errorf("Cacher.Put() extra function calls inconsistent")
	}
	if !invalidator.checkMetadata(expData) {
		t.Errorf(
			"Cacher.Put() metadata is incorrect, expected:\n%s\ngot:\n%s\n",
			expData, invalidator.getMetadata(),
		)
	}
	if !IsValueNotPresentError(err) {
		if err == nil {
			t.Errorf(
				"Cacher.Get() should have returned a ValueNotPresentError for a missing key",
			)
		} else {
			t.Errorf(
				"Cacher.Get() returned '%s' instead of a ValueNotPresentError", err,
			)
		}
	}
	if val != nil {
		t.Errorf("Cacher.Get() shouuld have returned nil for a missing value")
	}
}

func pString(s string) *string {
	return &s
}

type dummyHandler map[string]interface{}

func (d dummyHandler) Put(key string, value interface{}) error {
	d[key] = value
	return nil
}

func (d dummyHandler) Get(key string) (interface{}, error) {
	var val interface{}
	var ok bool
	if val, ok = d[key]; !ok {
		return nil, ValueNotPresentError{Key: key}
	}
	return val, nil
}

func (d dummyHandler) Clear() error {
	for k, _ := range d {
		delete(d, k)
	}
	return nil
}

func (d dummyHandler) Remove(key string) error {
	if _, ok := d[key]; !ok {
		return ValueNotPresentError{Key: key}
	}
	delete(d, key)
	return nil
}

func (d dummyHandler) Range(cb func(string, interface{}) bool) {
	for k, v := range d {
		if !cb(k, v) {
			return
		}
	}
}

func ExampleNewInMemoryDataHandler() {
	// InMemoryDataHandler is the default.
	// could also do: NewCache(NewInMemoryDataHandler(), nil)
	myCache := NewCache(nil, nil)
	newItem, _ := myCache.Get("foo", 42)
	found, _ := myCache.Get("foo")
	if newItem.(int) != found.(int) {
		fmt.Println("not equal")
	} else {
		fmt.Println("equal")
	}
	// Output: equal
}
