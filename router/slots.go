package router

import (
	"errors"

	"bitbucket.org/nedp/command"
)

type Slots interface {
	// Add finds a free slot and assigns it to the specified command.
	//
	// Returns
	// the index of the slot assigned to the command.
	Add(c command.Interface) (int, error)

	// Run is a wrapper for Run on the command in slot i.
	// i must be positive, slot i must be assigned to a command, and
	// the command must not already be running.
	//
	// Returns
	// an error if there is a problem running the command; and
	// command.Run's result if the command runs successfully.
	Run(i int, outCh chan<- string) (bool, error)

	// Free unassigns the slot with index i.
	// i must be positive, slot i must be assigned to a command, and
	// the command must not be running.
	//
	// Returns
	// nil if successful;
	// ErrNotAssigned if slot i is not assigned; and
	// ErrStillRunning if the command in slot i is running.
	Free(i int) error
}

type slots struct {
	commands []command.Interface
	nUsed int
	iSlot int
	maxNSlots int
}

var (
	ErrNotAssigned = errors.New("slots: tried to use the assignee of an unasigned slot")
	ErrStillRunning = errors.New("slots: tried to free a slot with a still-running command")
	ErrNoFreeSlots = errors.New("slots: tried to assign a slot when none are free")
)

const growthRate = 2
const sparsityFactor = 2

func NewSlots(nSlots int, maxNSlots int) Slots {
	// Preconditions
		if nSlots < 0 {
			panic("router.NewSlots: nSlots out of range")
		}

	return &slots{
		commands: make([]command.Interface, nSlots),
		nUsed: 0,
		iSlot: 0,
		maxNSlots: maxNSlots,
	}
}

func (s *slots) Add(c command.Interface) (int, error) {
	if s.nUsed == s.maxNSlots {
		return 0, ErrNoFreeSlots
	}
	// Add new `nil` slots if it's getting crowded, up to the maximum.
	if (len(s.commands) < s.maxNSlots) && (s.nUsed*sparsityFactor >= len(s.commands)) {
		targetNSlots := growthRate * len(s.commands)
		if targetNSlots == len(s.commands) {
			targetNSlots += 1
		}
		if targetNSlots > s.maxNSlots {
			targetNSlots = s.maxNSlots
		}
		nNewSlots := targetNSlots - len(s.commands)

		newSlots := make([]command.Interface, nNewSlots)

		s.iSlot = len(s.commands)
		s.commands = append(s.commands, newSlots...)
	}

	// Find a free slot
	i := s.iSlot
	for s.commands[i] != nil {
		if i + 1 == len(s.commands) {
			i = 0
		}
		i += 1
	}
	// Add the command
	s.commands[i] = c
	s.nUsed += 1

	// Skip the next slot in a future search to maintain sparsity,
	// reducing the expected number of checks.
	s.iSlot = i + sparsityFactor
	for s.iSlot >= len(s.commands) {
		s.iSlot -= len(s.commands)
	}

	return i, nil
}

func (s *slots) Run(i int, outCh chan<- string) (bool, error) {
	// Preconditions
		if i < 0 {
			panic("Router.*slots.Run: i out of range")
		}
		if s.commands[i] == nil {
			return false, ErrNotAssigned
		}
		if s.commands[i].IsRunning() {
			return false, ErrStillRunning
		}

	return s.commands[i].Run(outCh), nil
}

func (s *slots) Free(i int) error {
	// Preconditions
		if i < 0 {
			panic("Router.*slots.Free: i out of range")
		}
		if s.commands[i] == nil {
			return ErrNotAssigned
		}
		if s.commands[i].IsRunning() {
			return ErrStillRunning
		}

	s.commands[i] = nil
	s.nUsed -= 1
	return nil
}
