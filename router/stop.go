package router

import (
	"fmt"
	"log"
	s "github.com/nedp/command/sequence"
)

type StopParams struct {
	ID int
	Router *Router
}

func (StopParams) IsParams() {} // Marker

// Need to use a closure to capture the router.
func NewStopParams(r *Router) func() Params {
	return func() Params {
		p := new(StopParams)
		p.Router = r
		return p
	}
}

func NewStopSequence(routeParams Params) s.RunAller {
	id := routeParams.(*StopParams).ID
	r := routeParams.(*StopParams).Router
	outCh := make(chan string)

	// Get the command, reporting failure if it's not there.
	sl := <-r.slots
	cmd, err := sl.Command(id)
	r.slots <- sl
	if err != nil {
		msg := fmt.Sprintf("Couldn't find the command in slot %d: %s", id, err.Error())
		return s.FirstJust(sendFailure(outCh, msg)).End(outCh)
	}

	// Send the command's name and slot.
	name := cmd.Name()
	builder := s.FirstJust(sendIntro(outCh, "Stopping", id, name))

	// Stop the command.
	// If a failure occurs, report it, close the channel, and end the sequence.
	err = cmd.Stop()
	if err != nil {
		builder = builder.ThenJust(sendFailure(outCh, err.Error()))
		log.Printf("failed to stop command `%s` in slot %d:\n\t%s", name, id, err.Error())
		return builder.End(outCh)
	}

	// Report success
	builder = builder.ThenJust(sendString(outCh, "Command is now stopped"))

	// Close the channel
	builder = builder.ThenJust(closeCh(outCh))

	log.Printf("stopping command `%s` in slot %d", name, id)
	return builder.End(outCh)
}
