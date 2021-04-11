package cache

import (
	"fmt"
	"testing"
	"time"
)

var testCases = []*struct {
	name   string
	tData  *Metadata
	durStr string
	exp    bool
}{
	//huh?
	{
		name: "create",
		tData: &Metadata{
			Created: time.Now().Unix(),
		},
		durStr: "5s",
		exp:    true,
	},
	{
		name: "access",
		tData: &Metadata{
			Accessed: time.Now().Unix(),
		},
		durStr: "5s",
		exp:    true,
	},
	{
		name: "modified",
		tData: &Metadata{
			Modified: time.Now().Unix(),
		},
		durStr: "5s",
		exp:    true,
	},
	{
		name: "create expired",
		tData: &Metadata{
			Created: time.Now().Unix() - 5,
		},
		durStr: "1s",
		exp:    false,
	},
	{
		name: "access expired",
		tData: &Metadata{
			Accessed: time.Now().Unix() - 5,
		},
		durStr: "1s",
		exp:    false,
	},
	{
		name: "modified expired",
		tData: &Metadata{
			Modified: time.Now().Unix() - 5,
		},
		durStr: "1s",
		exp:    false,
	},
	{
		name: "complex not expired",
		tData: &Metadata{
			Created:  time.Now().Unix() - 5,
			Modified: time.Now().Unix() - 1,
			Accessed: time.Now().Unix() - 2,
		},
		durStr: "10s",
		exp:    true,
	},
	{
		name: "complex expired",
		tData: &Metadata{
			Created:  time.Now().Unix() - 5,
			Modified: time.Now().Unix() - 10,
			Accessed: time.Now().Unix() - 20,
		},
		durStr: "1s",
		exp:    false,
	},
}

func TestTimedInvalidator(t *testing.T) {
	for i, tCase := range testCases {
		dur, _ := time.ParseDuration(tCase.durStr)
		inv := NewTimedInvalidator(dur)
		if res := inv.IsValid(tCase.tData); res != tCase.exp {
			t.Errorf(
				"test '%s', index %d -- expected %t, got %t, ptr: %p",
				tCase.name, i, tCase.exp, res, inv,
			)
		}
	}
}

func ExampleNewTimedInvalidator() {
	lifetime, _ := time.ParseDuration(".5s")
	myCache := NewCache(nil, NewTimedInvalidator(lifetime))
	myCache.Put("foo", "a string")
	oneSec, _ := time.ParseDuration("1s")
	// IsValid is run every half second, make sure enough time has passed.
	time.Sleep(2 * oneSec)
	item, err := myCache.Get("foo")
	if item != nil || !IsValueNotPresentError(err) {
		fmt.Println("won't see this")
	} else {
		fmt.Println("expired")
	}
	// Output: expired
}
