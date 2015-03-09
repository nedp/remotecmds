package router

import (
	"fmt"
)

func sendString(outCh chan<- string, str string) func() error {
	return func() error {
		outCh <- str
		return nil
	}
}

func sendCommandName(outCh chan<- string, id int, name string) func() error {
	return sendString(outCh, fmt.Sprintf("Command (%s) in slot %d:", name, id))
}

func sendOutputLine(outCh chan<- string, iLine int, line string) func() error {
	return sendString(outCh, fmt.Sprintf("Output %4d: %s", iLine, line))
}

func sendFailure(outCh chan<- string, msg string) func() error {
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
