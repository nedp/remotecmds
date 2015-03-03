package router

import (
	"fmt"
	"log"

	"bitbucket.org/nedp/command"
	"bitbucket.org/nedp/command/sequence"
)

const defaultRoutesCapacity = 6

type Router struct {
	routes map[string]func() Route
	slots chan *slots
}

type Interface interface {
	RouteFor(string) (Route, error)
	SequenceFor(string) (sequence.RunAller, error)
	AddRoute(name string, newSeq func(Params) sequence.RunAller, newPs func() Params)
	OutputFor(req string) (<-chan string, error)
}

func (r Router) RouteFor(request string) (Route, error) {
	return RouteFor(request, r.routes)
}

func (r Router) SequenceFor(request string) (sequence.RunAller, error) {
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

// New creates and returns a new Router,
// with a command pool of the specified size.
//
// nSlots is only used as an initial number of slots to
// make allocation more efficient.
// More slots will be allocated if needed, to a maximum of
// maxNSlots.
func New(nSlots int, maxNSlots int) Interface {
	cr := &Router{
		make(map[string]func() Route, defaultRoutesCapacity),
		make(chan *slots, 1),
	}
	cr.slots <- newSlots(nSlots, maxNSlots)
	return cr
}

// OutputFor routes a request, generating its sequence,
// then creating and running a new command for it.
//
// If there are currently no free slots for commands,
// more will be created.
//
// The request is routed using the routes already registered
// with the Router.
// A new output channel is created, passed to the command,
// and returned to the caller.
//
// Returns
// (the output channel, nil) if the routing succeeds; and
// (nil, an error) if the routing fails.
func (cr *Router) OutputFor(req string) (<-chan string, error) {
	// Resolve the route
	rt, err := cr.RouteFor(req)
	if err != nil {
		return nil, err
	}
	seq := rt.Sequence()

	// Make the command
	if err != nil {
		return nil, fmt.Errorf("couldn't route the request: %s", err.Error())
	}
	s := <-cr.slots
	defer func() { cr.slots <- s }()

	cmd := command.New(seq)
	iSlot, err := s.Add(cmd)
	if err != nil {
		return nil, fmt.Errorf("couldn't add a new command: %s", err.Error())
	}
	outCh := make(chan string, 1)

	// Run the new command
	// TODO error handling other than printing logs and crashing.

	log.Printf("running command (%s) in slot %d", rt.Name, iSlot)
	outCh <- fmt.Sprintf("%s running in slot %d", rt.Name, iSlot)

	go func(sCh chan *slots, cmd command.Interface, iSlot int, outCh chan<- string, name string) {
		ok := cmd.Run(outCh)

		if err != nil {
			log.Fatalf(
				"cmdrouter.OutputFor: slot reported an error running a new command:\nt%s",
				err.Error(),
			)
		}
		if ok {
			log.Printf("command (%s) in slot %d completed successfully", name, iSlot)
		} else {
			log.Printf("command (%s) in slot %d failed", name, iSlot)
		}
		s := <-sCh
		err = s.Free(iSlot)
		sCh <- s
		if err != nil {
			log.Printf("couldn't free slot %d: %s", iSlot, err.Error())
		}
	}(cr.slots, s.commands[iSlot], iSlot, outCh, rt.Name)

	return outCh, nil
}
