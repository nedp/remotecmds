/*
Package say provides an interruptible speaking service.
The server will 'say' a user supplied quote.

It breaks the quote into phrases, using punctuation as delimiters,
running the espeak command for each phrase.
This allows the sequence to be terminated at any point
between phrases, but introduces a short pause at punctuation.
*/
package say

import (
	"errors"
	"io"
	"os/exec"
	"strings"

	s "bitbucket.org/nedp/command/sequence"
	"bitbucket.org/nedp/remotecmds/router"
)

type Params struct {
	Quote string
}
func (Params) IsParams() {} // Marker
func NewParams() router.Params {
	return new(Params)
}

// Creates a new sequence (bitbucket.org/nedp/command/sequence)
// for saying `quote`.
//
// Returns
// the created sequence.
func NewSequence(routeParams router.Params) s.RunAller {
	p := params{
		routeParams.(*Params).Quote,
		make(chan *exec.Cmd, 1),
		make(chan io.WriteCloser, 1),
	}
	println("quote:", p.quote)

	phrases := phrasesIn(p.quote)

	builder := s.SequenceOf(
		// Prepare the first espeak instance and pipe.
		s.PhaseOf(prepareEspeak(p)).
			// Send through the first pipe.
			AndJust(sendPhrase(p, phrases[0])),
	)

	for i := 0; i+1 < len(phrases); i += 1 {
		builder = builder.Then(
			// Run the current espeak instance.
			s.PhaseOf(runEspeak(p)).
				// Prepare the next espeak instance and pipe.
				AndJust(prepareEspeak(p)).
				// Send through the next pipe.
				AndJust(sendPhrase(p, phrases[i+1])),
		)
	}

	// Run the last espeak instance.
	builder = builder.ThenJust(runEspeak(p))

	outCh := make(chan string)
	close(outCh)

	return builder.End(outCh)
}

func phrasesIn(quote string) []string {
	const sep = ",.;:?!"
	phrases := make([]string, 0, len(quote)/2)
	// Break the quote into phrases on seperator characters.
	for i := strings.IndexAny(quote, sep); i != -1; i = strings.IndexAny(quote, sep) {
		// If there's a separator character at the end of the quote,
		// add the rest of the quote as a single phrase
		// and break from the loop early.
		if i+1 == len(quote) {
			break
		}
		// Add the next phrase.
		phrases = append(phrases, quote[:i+1])
		quote = quote[i+1:]
	}
	// Add the part of the quote after the last punctuation mark.
	return append(phrases, quote[:])
}

type params struct {
	quote string

	cmdCh chan *exec.Cmd
	pipeCh chan io.WriteCloser
}

func prepareEspeak(p params) func() error {
	return func() error {
		cmd := exec.Command("espeak")
		if pipe, err := cmd.StdinPipe(); err != nil {
			return err
		} else {
			p.pipeCh <- pipe
		}
		p.cmdCh <- cmd
		return nil
	}
}

func sendPhrase(p params, phrase string) func() error {
	return func() error {
		pipe := <-p.pipeCh
		for n, err := pipe.Write([]byte(phrase)); n < len(phrase);
				n, err = pipe.Write([]byte(phrase)) {
			if n <= 0 {
				return errors.New("Failed to write to the pipe")
			}
			if err != nil {
				return err
			}
			phrase = phrase[n:]
		}
		pipe.Close()
		return nil
	}
}

func runEspeak(p params) func() error {
	return func() error {
		cmd := <-p.cmdCh
		return cmd.Run()
	}
}

func splitFunc(c rune) bool {
	switch c {
	case ',', '.', ';', ':', '?', '!':
		return true
	default:
		return false
	}
}
