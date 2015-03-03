package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const longMax = 10000
const shortMax = 100

type commandMock struct {
	mock.Mock
}

func (c *commandMock) Output() []string {
	return c.Called().Get(0).([]string)
}

func (c *commandMock) Run(chan<- string) bool {
	return c.Called().Bool(0)
}

func (c *commandMock) Pause() (bool, error) {
	args := c.Called()
	return args.Bool(0), args.Error(1)
}

func (c *commandMock) Cont() (bool, error) {
	args := c.Called()
	return args.Bool(0), args.Error(1)
}

func (c *commandMock) Stop() error {
	return c.Called().Error(0)
}

func (c *commandMock) IsRunning() bool {
	return c.Called().Bool(0)
}

func TestNewSlots(t *testing.T) {
	for i := 0; i < longMax; i += 1 {
		s := NewSlots(i, i).(*slots)
		assert.Equal(t, i, len(s.commands), "newSlots produced the wrong number of slots")
		assert.Equal(t, 0, s.nUsed, "newSlots produced non-empty slots")
		assert.Equal(t, 0, s.iSlot, "newSlots produced non-zeroed iSlot object")
	}
}

func TestFree(t *testing.T) {
	for iCm := 0; iCm < shortMax; iCm += 1 {
		s := NewSlots(shortMax, shortMax).(*slots)
		cm := new(commandMock)
		cm.On("IsRunning").Return(false) // Want to free it immediately

		s.commands[iCm] = cm
		for i := 0; i < iCm; i += 1 {
			assert.NotNil(t, s.Free(i), "Freed an unused slot, but shouldn't be able to.")
		}
		assert.Nil(t, s.Free(iCm), "Couldn't free a used slot.")
		for i := iCm + 1; i < shortMax; i += 1 {
			assert.NotNil(t, s.Free(i), "Freed an unused slot, but shouldn't be able to.")
		}
	}
}

func TestFreeStillRunning(t *testing.T) {
	cm := new(commandMock)
	cm.On("IsRunning").Return(true) // Want to free it immediately
	s := NewSlots(1, 1).(*slots)
	s.commands[0] = cm
	assert.Equal(t, ErrStillRunning, s.Free(0),
		"Freed a still running slot, didn't get expected error.")
}

func TestFreeAlreadyFree(t *testing.T) {
	s := NewSlots(1, 1).(*slots)
	assert.Equal(t, ErrNotAssigned, s.Free(0),
		"Freed an unassigned slot, didn't get expected error.")
}

func TestUseIncrease(t *testing.T) {
	cm := new(commandMock)
	s := NewSlots(0, 1).(*slots)
	_, err := s.Add(cm)
	assert.Nil(t, err, "Error extending number of slots")
}

func TestUseFull(t *testing.T) {
	cm := new(commandMock)
	s := NewSlots(0, 0).(*slots)
	_, err := s.Add(cm)
	assert.Equal(t, ErrNoFreeSlots, err,
		"Didn't get expected error adding to full slots.")
}
