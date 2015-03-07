package utc

import (
	"time"

	s "bitbucket.org/nedp/command/sequence"
	"bitbucket.org/nedp/remotecmds/router"
)

func NewSequence(_ router.Params) s.RunAller {
	outCh := make(chan string)

	b := s.FirstJust(SendUTC(outCh))
	return b.End(outCh)
}

func SendUTC(ch chan<- string) func() error {
	return func() error {
		ch <- time.Now().UTC().String()
		close(ch)
		return nil
	}
}
