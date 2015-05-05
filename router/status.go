package router

import (
	"fmt"
	"log"

	s "github.com/nedp/command/sequence"
)

type StatusParams struct {
	ID int
	Router *Router
}

func (StatusParams) IsParams() {} // Marker

// Need to use a closure to capture the router.
func NewStatusParams(r *Router) func() Params {
	return func() Params {
		p := new(StatusParams)
		p.Router = r
		return p
	}
}

func NewStatusSequence(routeParams Params) s.RunAller {
	id := routeParams.(*StatusParams).ID
	r := routeParams.(*StatusParams).Router
	outCh := make(chan string)

	// Get information about target command.
	// Fail this command if we can't find the target command.
	sl := <-r.slots
	cmd, err := sl.Command(id)
	r.slots <- sl
	if err != nil {
		msg := fmt.Sprintf("Couldn't find the command in slot %d: %s", id, err.Error())
		return s.FirstJust(sendFailure(outCh, msg)).End(outCh)
	}
	state := cmd.State()
	name := cmd.Name()

	// Send the command name
	builder := s.FirstJust(sendIntro(outCh, "Status of", id, name))

	// If the command is running, report it.
	if state.IsRunning {
		builder = builder.ThenJust(sendString(outCh, "Command is still Running"))
	// Otherwise report whether the command has been stopped, or just finished running.
	} else {
		if state.HasStopped {
			builder = builder.ThenJust(sendString(outCh, "Command has been stopped"))
		} else {
			builder = builder.ThenJust(sendString(outCh, "Command has finished successfully"))
		}
	}

	// If the command is paused, report it
	if state.IsPaused {
		builder = builder.ThenJust(sendString(outCh, "Command is currently paused"))
	}

	// Send each line in the command's prior output
	for i, line := range state.Output {
		builder = builder.ThenJust(sendOutputLine(outCh, i, line))
	}

	// Close the channel
	builder = builder.ThenJust(closeCh(outCh))

	log.Printf("reporting status of command `%s` in slot %d", name, id)
	return builder.End(outCh)
}
