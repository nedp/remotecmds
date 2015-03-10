package router

import (
	"errors"
	"fmt"
)

func sendString(outCh chan<- string, str string) func() error {
	return func() error {
		outCh <- str
		return nil
	}
}

func sendIntro(outCh chan<- string, prefix string, id int, name string) func() error {
	return sendString(outCh, fmt.Sprintf("%s command `%s` in slot %d:", prefix, name, id))
}

func sendOutputLine(outCh chan<- string, iLine int, line string) func() error {
	return sendString(outCh, fmt.Sprintf("Output %4d: %s", iLine, line))
}

func sendFailure(outCh chan<- string, msg string) func() error {
	return func() error {
		newMsg := fmt.Sprintf("Error: %s", msg)
		outCh <- newMsg
		close(outCh)
		return errors.New(newMsg)
	}
}

func closeCh(ch chan<- string) func() error {
	return func() error {
		close(ch)
		return nil
	}
}
