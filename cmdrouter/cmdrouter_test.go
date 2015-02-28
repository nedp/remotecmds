package cmdrouter

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type commandMock struct {
	outCh chan<- string
}

func (ra *runAllerMock) RunAll
