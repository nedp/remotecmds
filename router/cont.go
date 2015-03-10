package router

import (
	"fmt"
	"log"
	s "bitbucket.org/nedp/command/sequence"
)

type ContParams struct {
	ID int
	Router *Router
}

func (ContParams) IsParams() {} // Marker

// Need to use a closure to capture the router.
func NewContParams(r *Router) func() Params {
	return func() Params {
		p := new(ContParams)
		p.Router = r
		return p
	}
}

func NewContSequence(routeParams Params) s.RunAller {
	id := routeParams.(*ContParams).ID
	r := routeParams.(*ContParams).Router
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
	builder := s.FirstJust(sendIntro(outCh, "Continuing", id, name))

	// Continue the command.
	// If a failure occurs, report it, close the channel, and end the sequence.
	wasAlreadyContinuing, err := cmd.Cont()
	if err != nil {
		builder = builder.ThenJust(sendFailure(outCh, err.Error()))
		log.Printf("failed to continue command `%s` in slot %d", name, id)
		return builder.End(outCh)
	}

	// Report whether the command was already continuing, or if we contued it.
	if wasAlreadyContinuing {
		builder = builder.ThenJust(sendString(outCh, "Command was already continuing"))
	} else {
		builder = builder.ThenJust(sendString(outCh, "Command is now continuing"))
	}

	// Close the channel
	builder = builder.ThenJust(closeCh(outCh))

	log.Printf("continuing command `%s` in slot %d", name, id)
	return builder.End(outCh)
}
