package say

import (
	"testing"
	"time"

	"bitbucket.org/nedp/command"
)

const shortDuration = time.Second

func TestCommand(t *testing.T) {
	command.New(NewSequence(Params{"This is the Integration test! Will it pass?"})).Run()
}

func TestKill(t *testing.T) {
	c := command.New(NewSequence(Params{"This should be cut off now; before it finishes."}))
	go c.Run()
	go func() {
		time.Sleep(shortDuration)
		c.Kill()
	}()
	<-c.WhenTerminated()
}
