package say

import (
	"io"
	"os/exec"
	"testing"

	//"bitbucket.org/nedp/command/sequence"
)

func TestUnits(t *testing.T) {
	p := params{
		"Hello, world!",
		make(chan *exec.Cmd, 1),
		make(chan io.WriteCloser, 1),
	}
	phrases := phrasesIn(p.quote)

	go sendPhrase(p, phrases[0])()
	prepareEspeak(p)()
	for i := 0; i+1 < len(phrases); i += 1 {
		go prepareEspeak(p)()
		go sendPhrase(p, phrases[i+1])()
		runEspeak(p)()
	}
	runEspeak(p)()
}

func TestUnitsFail(t *testing.T) {
	p := params{
		"Hello, world!",
		make(chan *exec.Cmd, 1),
		make(chan io.WriteCloser, 1),
	}
	phrases := phrasesIn(p.quote)

	go sendPhrase(p, phrases[0])()
	prepareEspeak(p)()
	for i := 0; i+1 < len(phrases); i += 1 {
		go prepareEspeak(p)()
		go sendPhrase(p, phrases[i+1])()
		runEspeak(p)()
		return // Failure
	}
	runEspeak(p)()
}
