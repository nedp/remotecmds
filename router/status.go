package router

import (
	"fmt"

	s "bitbucket.org/nedp/command/sequence"
)

type StatusParams struct {
	ID int
	Router *Router
}

func (StatusParams) IsParams() {} // Marker

// Need to use a closure to capture the router.
func NewParamsStatus(r *Router) func() Params {
	return func() Params {
		p := new(StatusParams)
		p.Router = r
		return p
	}
}

func NewSequenceStatus(routeParams Params) s.RunAller {
	id := routeParams.(*StatusParams).ID
	r := routeParams.(*StatusParams).Router
	outCh := make(chan string)

	// Get information about target command.
	// Fail this command if we can't find the target command.
	sl := <-r.slots
	cmd, err := sl.Command(id)
	r.slots <- sl
	if err != nil {
		msg := fmt.Sprintf("ERROR: Couldn't find the command: %s", err.Error())
		return s.FirstJust(fail(outCh, msg)).End(outCh)
	}
	state := cmd.State()
	name := cmd.Name()

	// Send the command name
	builder := s.FirstJust(sendName(outCh, id, name))

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
		builder = builder.ThenJust(sendOut(outCh, i, line))
	}

	// Close the channel
	builder = builder.ThenJust(closeCh(outCh))

	return builder.End(outCh)
}

func sendString(outCh chan<- string, str string) func() error {
	return func() error {
		outCh <- str
		return nil
	}
}

func sendName(outCh chan<- string, id int, name string) func() error {
	return sendString(outCh, fmt.Sprintf("Command (%s) in slot %d:", name, id))
}

func sendOut(outCh chan<- string, iLine int, line string) func() error {
	return sendString(outCh, fmt.Sprintf("Output line %4d: %s", iLine, line))
}

func fail(outCh chan<- string, msg string) func() error {
	return func() error {
		outCh <- msg
		close(outCh)
		return fmt.Errorf("status command failed: %s", msg)
	}
}

func closeCh(ch chan<- string) func() error {
	return func() error {
		close(ch)
		return nil
	}
}
