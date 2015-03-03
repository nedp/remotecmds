package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testNSlots = 8
	testMaxSlots = 16
)

func TestRouteForMethod(t *testing.T) {
	r := &Router{testRoutes, make(chan *slots, 1)}
	r.slots <- newSlots(testNSlots, testMaxSlots)
	a, errA := RouteFor(testString, testRoutes)
	b, errB := r.RouteFor(testString)

	assert.Equal(t, a.Name, b.Name, "RouteFor method didn't match RouteFor function.")
	assert.Equal(t, a.NewSequence, b.NewSequence, "RouteFor method didn't match RouteFor function.")
	assert.Equal(t, *a.Params.(*ParamsTest), *b.Params.(*ParamsTest), "RouteFor method didn't match RouteFor function.")
	assert.Equal(t, errA, errB, "RouteFor method didn't match RouteFor function.")
}

func TestSequenceForMethod(t *testing.T) {
	r := &Router{testRoutes, make(chan *slots, 1)}
	r.slots <- newSlots(testNSlots, testMaxSlots)

	a, errA := SequenceFor(testString, testRoutes)
	b, errB := r.SequenceFor(testString)

	assert.Equal(t, a, b, "SequenceFor method didn't match SequenceFor function.")
	assert.Equal(t, errA, errB, "SequenceFor method didn't match SequenceFor function.")
}
