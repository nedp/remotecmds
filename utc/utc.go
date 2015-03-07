package utc

import (
	"time"

	s "bitbucket.org/nedp/command/sequence"
	"bitbucket.org/nedp/remotecmds/router"
)

type Params struct {
	NeedLocalTime bool `toml:"local"`
}

func (Params) IsParams() {} // Marker

func NewParams() router.Params {
	return new(Params)
}

func NewSequence(routeParams router.Params) s.RunAller {
	outCh := make(chan string)
	var b s.SequenceBuilder
	if routeParams.(*Params).NeedLocalTime {
		b = s.FirstJust(SendLocal(outCh))
	} else {
		b = s.FirstJust(SendUTC(outCh))
	}
	return b.End(outCh)
}

func SendUTC(ch chan<- string) func() error {
	return func() error {
		ch <- time.Now().UTC().String()
		close(ch)
		return nil
	}
}

func SendLocal(ch chan<- string) func() error {
	return func() error {
		ch <- time.Now().String()
		close(ch)
		return nil
	}
}
