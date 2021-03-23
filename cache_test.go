package cache

import (
	"testing"
)

func TestCacher(t *testing.T) {
	t.Run("type=inMemoryCache", testInMemoryStringCache)
}

type dummyStruct struct {
	number   int
	words    string
	contents *dummyStruct
}

func mkString(str string) *string {
	return &str
}

//TODO, need to test Get(sting, ...string)
//perhaps add a Empty() bool?
//test reset and clear
var inMemoryStringCacheTests = []*struct {
	name   string
	keys   []string
	values []*string
}{
	{
		"single key/value",
		[]string{"foo"},
		[]*string{mkString("bar")},
	},
	{
		"multiple values",
		[]string{"foo", "bar", "baz"},
		[]*string{mkString("blarg"), mkString("biff"), mkString("bar")},
	},
}

func testInMemoryStringCache(t *testing.T) {
	for _, tCase := range inMemoryStringCacheTests {
		myCache := NewInMemoryStringCacher()
		for i, k := range tCase.keys {
			err := myCache.Put(k, *tCase.values[i])
			if err != nil {
				t.Errorf(
					"should not have gotten an error for '%s' test, key: '%s', value '%s'",
					tCase.name, k, *tCase.values[i],
				)
				break
			}
			cVal, err := myCache.Get(k)
			if err != nil {
				t.Errorf(
					"should not have gotten an error for '%s' test, key: '%s', value '%s'",
					tCase.name, k, *tCase.values[i],
				)
				break
			}
			if cVal != *tCase.values[i] {
				t.Errorf(
					"%s test: expected value '%s' does not match returned value, '%s', for Get()",
					tCase.name, *tCase.values[i], cVal,
				)
			}
		}
	}
}
