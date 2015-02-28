package router

import (
	"bitbucket.org/nedp/command/sequence"
)

const defaultRoutesCapacity = 6

type RunAllerMaker func(Params) (sequence.RunAller, <-chan string)

type Interface interface {
	RouteFor([]byte) (Route, error)
	SequenceFor([]byte) (sequence.RunAller, <-chan string, error)
	AddRoute(name string, newSequence RunAllerMaker, newParams func() Params)
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

func (r Router) SequenceFor(request []byte) (sequence.RunAller, <-chan string, error) {
	return SequenceFor(request, r.routes)
}

func (r *Router) AddRoute(name string, newSequence RunAllerMaker, newParams func() Params) {
	r.routes[name] = func() Route {
		return Route{
			name,
			newSequence,
			newParams(), // Not called until after a route is retrieved from the map!
		}
	}
}
