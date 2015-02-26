package router

import (
	"bitbucket.org/nedp/command/sequence"
)

type Router interface {
	SequenceFor([]byte) (sequence.RunAller, error)
}

type router struct {
	//TODO
}

func (rt *router) SequenceFor(request []byte) (sequence.RunAller, error) {
	return nil, nil // TODO
}
