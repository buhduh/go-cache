package invalidator

import (
	"testing"
	"time"
)

var testCases = []*struct {
	name   string
	tData  *Metadata
	durStr string
	exp    bool
}{
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
				"test '%s', index %d -- expected %t, got %t",
				tCase.name, i, tCase.exp, res,
			)
		}
		inv.Stop()
	}
}
