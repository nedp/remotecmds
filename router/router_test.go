package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testRouter = &Router{testRoutes}

func TestRouteForMethod(t *testing.T) {
	a, errA := RouteFor([]byte(testString), testRoutes)
	b, errB := testRouter.RouteFor([]byte(testString))

	assert.Equal(t, a.Name, b.Name, "RouteFor method didn't match RouteFor function.")
	assert.Equal(t, a.NewSequence, b.NewSequence, "RouteFor method didn't match RouteFor function.")
	assert.Equal(t, *a.Params.(*ParamsTest), *b.Params.(*ParamsTest), "RouteFor method didn't match RouteFor function.")
	assert.Equal(t, errA, errB, "RouteFor method didn't match RouteFor function.")
}

func TestSequenceForMethod(t *testing.T) {
	a, errA := SequenceFor([]byte(testString), testRoutes)
	b, errB := testRouter.SequenceFor([]byte(testString))

	assert.Equal(t, a, b, "SequenceFor method didn't match SequenceFor function.")
	assert.Equal(t, errA, errB, "SequenceFor method didn't match SequenceFor function.")
}
