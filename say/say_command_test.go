package say

import (
	"testing"
	"time"

	"bitbucket.org/nedp/command"
)

const shortDuration = time.Second

func TestCommand(t *testing.T) {
	seq, ch := NewSequence(Params{"This is the Integration test! Will it pass?"})
	command.New(seq, ch).Run(make(chan string))
}

func TestKill(t *testing.T) {
	seq, ch := NewSequence(Params{"This should be cut off now; before it finishes."})
	c := command.New(seq, ch)
	go c.Run(make(chan string))
	go func() {
		time.Sleep(shortDuration)
		c.Kill()
	}()
	<-c.WhenTerminated()
}
