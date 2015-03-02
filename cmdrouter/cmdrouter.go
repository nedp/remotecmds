package cmdrouter

import (
	"fmt"
	"log"

	"bitbucket.org/nedp/command"
	"bitbucket.org/nedp/remotecmds/router"
	"bitbucket.org/nedp/remotecmds/slots"
)

type CommandRouter struct {
	router.Interface
	slots chan slots.Interface
}

type Interface interface {
	router.Interface
	OutputFor(req string) (<-chan string, error)
}

// New creates and returns a new CommandRouter,
// with a command pool of the specified size.
//
// nSlots is only used as an initial number of slots to
// make allocation more efficient.
// More slots will be allocated if needed, to a maximum of
// maxNSlots.
func New(nSlots int, maxNSlots int) Interface {
	cr := &CommandRouter{
		router.New(),
		make(chan slots.Interface, 1),
	}
	cr.slots <- slots.New(nSlots, maxNSlots)
	return cr
}

// OutputFor routes a request, generating its sequence,
// then creating and running a new command for it.
//
// If there are currently no free slots for commands,
// more will be created.
//
// The request is routed using the routes already registered
// with the CommandRouter.
// A new output channel is created, passed to the command,
// and returned to the caller.
//
// Returns
// (the output channel, nil) if the routing succeeds; and
// (nil, an error) if the routing fails.
func (cr *CommandRouter) OutputFor(req string) (<-chan string, error) {
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

	cmd := command.New(seq)
	iSlot, err := s.Add(cmd)
	if err != nil {
		return nil, fmt.Errorf("couldn't add a new command: %s", err.Error())
	}
	outCh := make(chan string, 1)

	// Run the new command
	// TODO error handling other than printing logs and crashing.
	go func(iSlot int, outCh chan<- string, name string) {
		ok, err := s.Run(iSlot, outCh)
		if err != nil {
			log.Fatal("cmdrouter.(*CommandRouter).OutputFor: slot reported an error running a new command.")
		}
		if ok {
			log.Printf("(%s) command completed successfully", name)
		} else {
			log.Printf("(%s) command failed", name)
		}
	}(iSlot, outCh, rt.Name)

	outCh <- fmt.Sprintf("%s running in slot: %d", rt.Name, iSlot)
	cr.slots <- s

	return outCh, nil
}
