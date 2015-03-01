package router

import (
	"bitbucket.org/nedp/command/sequence"
)

const defaultRoutesCapacity = 6

type Interface interface {
	RouteFor([]byte) (Route, error)
	SequenceFor([]byte) (sequence.RunAller, error)
	AddRoute(name string, newSeq func(Params) sequence.RunAller, newPs func() Params)
}

type Router struct {
	routes map[string]func() Route
}

// Returns
// a new Router
func New() Interface {
	return &Router{make(map[string]func() Route, defaultRoutesCapacity)}
}

func (r Router) RouteFor(request []byte) (Route, error) {
	return RouteFor(request, r.routes)
}

func (r Router) SequenceFor(request []byte) (sequence.RunAller, error) {
	return SequenceFor(request, r.routes)
}

func (r *Router) AddRoute(name string, newSeq func(Params) sequence.RunAller,
		newPs func() Params,
) {
	r.routes[name] = func() Route {
		return Route{
			name,
			newSeq,
			newPs(), // Not called until after a route is retrieved from the map!
		}
	}
}
