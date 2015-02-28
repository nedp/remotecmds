package router

import (
	"testing"
	"time"

	"bitbucket.org/nedp/command/sequence"
	"bitbucket.org/nedp/command/status"

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

var testRoutes = map[string]func() Route {
	"test": func() Route {
		return Route{
			"test",
			func(_ Params) (sequence.RunAller, <-chan string) {
				return testDummySequence, testCh
			},
			new(ParamsTest),
		}
	},
}

type dummySequence struct {}
func (d *dummySequence) RunAll(s status.Interface) status.Interface {
	return s
}

func (d *dummySequence) IsRunning() bool {
	return false
}

var testDummySequence = new(dummySequence)
var testCh = make(<-chan string)

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
	assert.Equal(t, "Hello, World!", p.AA, "wrong Params.AA")
	assert.Equal(t, 5, p.BB, "wrong Params.BB")
}

func TestRouteForRequest(t *testing.T) {
	rt, err := RouteFor([]byte(testString), testRoutes)
	if err != nil {
		t.Fatal(err.Error())
	}
	p := rt.Params.(*ParamsTest)
	assert.Equal(t, "Hello, World!", p.AA, "wrong Params.AA")
	assert.Equal(t, 5, p.BB, "wrong Params.BB")

	assert.Equal(t, "test", rt.Name)
}

func TestSequenceFor(t *testing.T) {
	seq, ch, err := SequenceFor([]byte(testString), testRoutes)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, testDummySequence, seq, "SequenceFor didn't return the expected test dummy.")
	assert.Equal(t, testCh, ch, "SequenceFor didn't return the expected channel.")
}

