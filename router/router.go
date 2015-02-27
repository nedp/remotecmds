package router

import (
	"fmt"
	"bitbucket.org/nedp/command/sequence"
)

type Router interface {
	SequenceFor([]byte) (sequence.RunAller, error)
}

// Creates a new Router, using the routes in the supplied map.
//
// The map should map route names to their Route structs.
// Multiple aliases for the same route are allowed, but
// must each have their own Route struct with individualised
// `Name` fields.
//
// Basic verification is performed on the Routes:
// the map key must match the Route `Name`,
// and the Route's `Fn` and `Params` fields must not be nil.
// This has a runtime cost, but since making a new router
// should not be a frequent task, it is considered acceptable.
//
// Returns
// the new Router, `nil`     if verification passes;
// `nil`,           an error if verification fails.
func New(routes map[string]Route) (Router, error) {
	for key, r := range routes {
		if r.Name != key || r.Fn == nil || r.Params == nil {
			return nil, fmt.Errorf("invalid route detected (%s: %x)", key, r)
		}
	}
	return &router{routes}, nil
}

type router struct {
	routes map[string]Route
}

func (rt *router) SequenceFor(request []byte) (sequence.RunAller, error) {
	return nil, nil // TODO
}
