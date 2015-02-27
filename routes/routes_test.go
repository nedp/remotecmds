package routes

import (
	"testing"

	"bitbucket.org/nedp/remotecmds/say"
	"bitbucket.org/nedp/remotecmds/router"

	"github.com/stretchr/testify/assert"
)

const testString = `
	say:
	Quote = "This is the route test."
	`

func TestRoute(t *testing.T) {
	r, err := router.RouteFor([]byte(testString), Routes)

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "say", r.Name, "Route name didn't match")
	assert.Equal(t, say.New, r.Fn, "Route function didn't match")
	assert.Equal(t, "This is the route test.", r.Params.(*say.Params).Quote,
		"Route quote parameter didn't match")
}
