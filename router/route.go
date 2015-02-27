package router

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"

	"bitbucket.org/nedp/command/sequence"
)

const seperator = ":\n"

func SequenceFor(request []byte, routes map[string]func() Route) (sequence.RunAller, error) {
	rt, err := RouteFor(request, routes)
	if err != nil {
		return nil, err
	}
	return rt.Sequence(), nil
}

type Route struct {
	Name string
	NewSequence RunAllerMaker
	Params Params
}

func RouteFor(request []byte, routes map[string]func() Route) (Route, error) {
	// Strip leading whitespace
	requestString := strings.TrimSpace(string(request))
	// Parse the command name and TOML params table
	var args string
	var routeName string
	split := strings.SplitN(requestString, seperator, 2)
	if len(split) != 2 {
		err := fmt.Errorf("couldn't parse request \"%s\"", string(request))
		return Route{}, err
	}
	routeName = split[0]
	args = split[1]

	// Unmarshal the TOML
	rtFn, ok := routes[routeName]
	if !ok {
		return Route{}, fmt.Errorf("route name \"%s\" not recognised", routeName)
	}
	rt := rtFn()

	if _, err := toml.Decode(args, rt.Params); err != nil {
		return Route{}, fmt.Errorf("couldn't unmarshall arguments (%s)", err.Error())
	}

	return rt, nil
}

func (rt *Route) Sequence() sequence.RunAller {
	return rt.NewSequence(rt.Params)
}

type Params interface {
	IsParams()
}

