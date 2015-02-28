package cmdrouter

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type runAllerMock struct {
	mock.Mock
}

type commandMock struct {
	mock.Mock
}

func (c *commandMock) Run(outCh chan<- string) bool {
}
