package say

import (
	"testing"
	"time"

	"bitbucket.org/nedp/command"
)

const shortDuration = time.Second

func TestCommand(t *testing.T) {
	seq := NewSequence(Params{"This is the Integration test! Will it pass?"})
	command.New(seq).Run(make(chan string))
}

func TestStop(t *testing.T) {
	seq := NewSequence(Params{"This should be cut off now; before it finishes."})
	c := command.New(seq)
	go func() {
		time.Sleep(shortDuration)
		c.Stop()
	}()
	c.Run(make(chan string))
}
