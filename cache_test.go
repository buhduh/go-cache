package cache

import (
	"fmt"
	"testing"
)

func TestCacher(t *testing.T) {
	t.Run("simple_functionality", simpleTest)
}

type dummyFoo struct {
	foo int
}

type dummyBar struct {
	bar string
}

type dummyReader struct{}

func (d *dummyReader) Read([]byte) (int, error) {
	return 10, nil
}

var p0, p1, p2 int = 0, 1, 2
var s0, s1, s2 string = "0", "1", "2"

//TODO need to test the reader too
var simpleTests = []*struct {
	name       string
	valueTypes []string
	values     [][]interface{}
}{
	{
		name:       "primitives",
		valueTypes: []string{"int", "string"},
		values: [][]interface{}{
			[]interface{}{0, 1, 2},
			[]interface{}{"0", "1", "2"},
		},
	},
	{
		name:       "primitive pointers",
		valueTypes: []string{"*int", "*string"},
		values: [][]interface{}{
			[]interface{}{&p0, &p1, &p2},
			[]interface{}{&s0, &s1, &s2},
		},
	},
	{
		name:       "composite types",
		valueTypes: []string{"foo", "bar"},
		values: [][]interface{}{
			[]interface{}{dummyFoo{0}, dummyFoo{1}, dummyFoo{2}},
			[]interface{}{dummyBar{"0"}, dummyBar{"1"}, dummyBar{"2"}},
		},
	},
	{
		name:       "composite pointer types",
		valueTypes: []string{"*foo", "*bar"},
		values: [][]interface{}{
			[]interface{}{&dummyFoo{0}, &dummyFoo{1}, &dummyFoo{2}},
			[]interface{}{&dummyBar{"0"}, &dummyBar{"1"}, &dummyBar{"2"}},
		},
	},
}

func simpleTest(t *testing.T) {
	for _, tCase := range simpleTests {
		for j, vType := range tCase.valueTypes {
			if err := verifyValue(vType, tCase.values[j]); err != nil {
				t.Errorf(
					"%s test failed, type: '%s', error: '%s'",
					tCase.name, vType, err,
				)
			}
		}
	}
}

func verifyValue(vType string, data []interface{}) error {
	var ok bool
	switch vType {
	case "int":
		var toCheck int
		for i, v := range data {
			if toCheck, ok = v.(int); !ok {
				return fmt.Errorf("expected interface to be of type int")
			} else if i != toCheck {
				return fmt.Errorf("expected %d, got %d", i, toCheck)
			}
		}
	case "string":
		var toCheck string
		for i, v := range data {
			if toCheck, ok = v.(string); !ok {
				return fmt.Errorf("expected interface to be of type string")
			} else if fmt.Sprintf("%d", i) != toCheck {
				return fmt.Errorf("expected %d, got %s", i, toCheck)
			}
		}
	case "*int":
		var toCheck *int
		for i, v := range data {
			if toCheck, ok = v.(*int); !ok {
				return fmt.Errorf("expected interface to be of type *int")
			} else if i != *toCheck {
				return fmt.Errorf("expected %d, got %d", i, *toCheck)
			}
		}
	case "*string":
		var toCheck *string
		for i, v := range data {
			if toCheck, ok = v.(*string); !ok {
				return fmt.Errorf("expected interface to be of type *string")
			} else if fmt.Sprintf("%d", i) != *toCheck {
				return fmt.Errorf("expected %d, got %s", i, *toCheck)
			}
		}
	case "foo":
		var toCheck dummyFoo
		for i, v := range data {
			if toCheck, ok = v.(dummyFoo); !ok {
				return fmt.Errorf("expected interface to be of type dummyFoo")
			} else if i != toCheck.foo {
				return fmt.Errorf("expected %d, got %d", i, toCheck.foo)
			}
		}
	case "bar":
		var toCheck dummyBar
		for i, v := range data {
			if toCheck, ok = v.(dummyBar); !ok {
				return fmt.Errorf("expected interface to be of type dummyBar")
			} else if fmt.Sprintf("%d", i) != toCheck.bar {
				return fmt.Errorf("expected %d, got %s", i, toCheck.bar)
			}
		}
	case "*foo":
		var toCheck *dummyFoo
		for i, v := range data {
			if toCheck, ok = v.(*dummyFoo); !ok {
				return fmt.Errorf("expected interface to be of type *dummyFoo")
			} else if i != (*toCheck).foo {
				return fmt.Errorf("expected %d, got %d", i, toCheck.foo)
			}
		}
	case "*bar":
		var toCheck *dummyBar
		for i, v := range data {
			if toCheck, ok = v.(*dummyBar); !ok {
				return fmt.Errorf("expected interface to be of type *dummyBar")
			} else if fmt.Sprintf("%d", i) != (*toCheck).bar {
				return fmt.Errorf("expected %d, got %s", i, toCheck.bar)
			}
		}
	}
	return nil
}
