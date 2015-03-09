package router

import (
	"fmt"
	"log"
	s "bitbucket.org/nedp/command/sequence"
)

type PauseParams struct {
	ID int
	Router *Router
}

func (PauseParams) IsParams() {} // Marker

// Need to use a closure to capture the router.
func NewPauseParams(r *Router) func() Params {
	return func() Params {
		p := new(PauseParams)
		p.Router = r
		return p
	}
}

func NewPauseSequence(routeParams Params) s.RunAller {
	id := routeParams.(*PauseParams).ID
	r := routeParams.(*PauseParams).Router
	outCh := make(chan string)

	// Get the command, reporting failure if it's not there.
	sl := <-r.slots
	cmd, err := sl.Command(id)
	r.slots <- sl
	if err != nil {
		msg := fmt.Sprintf("ERROR: Couldn't find the command: %s", err.Error())
		return s.FirstJust(sendFailure(outCh, msg)).End(outCh)
	}

	// Send the command's name and slot.
	name := cmd.Name()
	builder := s.FirstJust(sendCommandName(outCh, id, name))

	// Pause the command.
	// If a failure occurs, report it, close the channel, and end the sequence.
	wasPaused, err := cmd.Pause()
	if err != nil {
		builder = builder.ThenJust(sendFailure(outCh, err.Error()))
		log.Printf("failed to pause command (%s) in slot %d", name, id)
		return builder.End(outCh)
	}

	// Report whether the command was already paused, or if we paused it.
	if wasPaused {
		builder = builder.ThenJust(sendString(outCh, "Command was already paused"))
	} else {
		builder = builder.ThenJust(sendString(outCh, "Command is now paused"))
	}

	// Close the channel
	builder = builder.ThenJust(closeCh(outCh))

	log.Printf("pausing command (%s) in slot %d", name, id)
	return builder.End(outCh)
}
