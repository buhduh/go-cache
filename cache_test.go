package cache

import (
//"testing"
)

//im tired, pretty sure this works...

/*
func TestCache(t *testing.T) {
	t.Run("Get", testGet)
}

var getTestCases = []*struct {
	name    string
	handler dummyHandler
	keys    []string
	vals    []interface{}
}{
	{
		name:    "one element",
		handler: map[string]interface{}{
      "key": "val",
    },
		keys:    []string{"foo"},
		vals:    []interface{}{"val"},
	},
}

func testGet(t *testing.T) {
  for i, tCase := range getTestCases {
    if len(tCase.keys) != len(tCase.vals) {
      t.Errorf("len(tCase.keys) != len(tCase.vals)")
    }
    tCache := NewCache(tCase.handler, &NopInvalidator{}, nil, nil, nil)
    for j, key := range tCase.keys {
      if res, err := tCache.Get(key);
    }
  }
}

type dummyHandler map[string]interface{}

func (d dummyHandler) Put(key string, data interface{}) error {
	d[key] = data
	return nil
}

//Must throw a ValueNotPresentError and be nil
func (d dummyHandler) Get(key string) (interface{}, error) {
	var toRet interface{}
	var ok bool
	if toRet, ok = d[key]; !ok || toRet == nil {
		return nil, ValueNotPresentError{Key: key}
	}
	return toRet, nil
}

func (d dummyHandler) Clear() error {
	for k, _ := range d {
		delete(d, k)
	}
	return nil
}

//Must throw a ValueNotPresentError
func (d dummyHandler) Remove(key string) error {
	var toRet interface{}
	var ok bool
	if toRet, ok = d[key]; !ok || toRet == nil {
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
*/
