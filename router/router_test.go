package router

import (
	"testing"
	"time"

	"bitbucket.org/nedp/command/sequence"

	"github.com/stretchr/testify/assert"
	"github.com/BurntSushi/toml"
)

type ParamsTest struct {
	AA string
	BB int
	CC time.Time
}

// Marker
func (ParamsTest) IsParams() {}

func Routes(name string) (route Route, ok bool) {
	ok = true
	switch name {
	case "test":
		route = Route{
			"test",
			func(_ Params) sequence.RunAller { return nil },
			new(ParamsTest),
		}
	default:
		ok = false
	}
	return
}

const testString = `test:
aa = "Hello, World!"
bb = 5
cc = 2015-02-25T16:11:00Z
`

const testArgs = `
aa = "Hello, World!"
bb = 5
cc = 2015-02-25T16:11:00Z
`

func TestParams(t *testing.T) {
	var p ParamsTest
	toml.Decode(testArgs, &p)
	assert.Equal(t, "Hello, World!", p.AA, "wrong Params.A")
	assert.Equal(t, 5, p.BB, "wrong Params.B")
}

func TestRouteForRequest(t *testing.T) {
	r, err := RouteForRequest([]byte(testString), Routes)
	if err != nil {
		t.Fatal(err.Error())
	}
	p := r.Params.(*ParamsTest)
	assert.Equal(t, "Hello, World!", p.AA, "wrong Params.A")
	assert.Equal(t, 5, p.BB, "wrong Params.B")

	assert.Equal(t, "test", r.Name)
}
