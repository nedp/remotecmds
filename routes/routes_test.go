package routes

import (
	"testing"

	"bitbucket.org/nedp/remotecmds/say"
	"bitbucket.org/nedp/remotecmds/router"

	"github.com/stretchr/testify/assert"
)

func TestSay(t *testing.T) {
	const testString = `
		say:
		Quote = "This is the route test."
		`
	rt := testHelper(t, "say", testString, say.NewSequence)
	assert.Equal(t, "This is the route test.", rt.Params.(*say.Params).Quote,
		"Route quote parameter didn't match")
}

// For each route, run this test helper, then check each of rt.Params' fields.
func testHelper(t *testing.T, name string, testString string, newSequence router.RunAllerMaker,
) router.Route {
	r := router.New()
	AddRoutesTo(r)

	rt, err := r.RouteFor([]byte(testString))

	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, name, rt.Name, "Route name didn't match")
	assert.Equal(t, newSequence,
		rt.NewSequence, "Route function didn't match")

	return rt
}
