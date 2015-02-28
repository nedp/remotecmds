package cmdrouter

import (
	"errors"
	"fmt"

	"bitbucket.org/nedp/command"
	"bitbucket.org/nedp/remotecmds/router"
)

type CommandRouter struct {
	router.Interface

	commandPool []command.Interface
	nCommands chan int
	iCommand chan int
}

type Interface interface {
	router.Interface
}

func New(poolSize int) Interface {
	cr := &CommandRouter{
		router.New(),
		make([]command.Interface, poolSize),
		make(chan int, 1),
		make(chan int, 1),
	}
	cr.nCommands <- 0
	cr.iCommand <- 0
	return cr
}

func (cr *CommandRouter) OutputFor(request []byte) (<-chan string, error) {
	// See if we have a free slot.
	nCommands := <-cr.nCommands
	defer func() { cr.nCommands <- nCommands }() // Late bound nCommands
	if nCommands >= len(cr.commandPool) {
		return nil, errors.New("Command pool is full.")
	}

	// Find a free slot.
	var iCommand int
	for iCommand = <-cr.iCommand; cr.commandPool[iCommand] != nil; iCommand += 1 {
		if iCommand + 1 == len(cr.commandPool) {
			iCommand = 0
		}
	}
	// Make the new command
	seq, seqOut, err := cr.SequenceFor(request)
	if err != nil {
		return nil, fmt.Errorf("couldn't route the request: %s", err.Error())
	}
	cr.commandPool[iCommand] = command.New(seq, seqOut)
	cr.iCommand <- iCommand + 1 // Don't need iCommand locked now

	// Run the new command
	outCh := make(chan string)
	cr.commandPool[iCommand].Run(outCh)

	return outCh, nil
}
