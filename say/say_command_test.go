package say

import (
	"testing"
	"time"

	"bitbucket.org/nedp/command"
)

const shortDuration = time.Second

func TestCommand(t *testing.T) {
	command.New(NewSequence(Params{"This is the Integration test! Will it pass?"}),
		make(chan string)).Run(make(chan string))
}

func TestKill(t *testing.T) {
	c := command.New(NewSequence(Params{"This should be cut off now; before it finishes."}),
		make(chan string))
	go c.Run(make(chan string))
	go func() {
		time.Sleep(shortDuration)
		c.Kill()
	}()
	<-c.WhenTerminated()
}
