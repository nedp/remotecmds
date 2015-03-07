package router

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"

	"bitbucket.org/nedp/command/sequence"
)

const seperator = ":"

func SequenceFor(req string, rts map[string]func() Route) (sequence.RunAller, error) {
	rt, err := RouteFor(req, rts)
	if err != nil {
		return nil, err
	}
	ra := rt.Sequence()
	return ra, nil
}

type Route struct {
	Name string
	NewSequence func(Params) sequence.RunAller
	Params Params
}

func RouteFor(request string, routes map[string]func() Route) (Route, error) {
	// Strip leading whitespace
	trimmedRequest := strings.TrimSpace(string(request))
	// Parse the command name and TOML params table
	var args string
	var routeName string
	split := strings.SplitN(trimmedRequest, seperator, 2)
	if len(split) != 2 {
		return Route{}, fmt.Errorf("couldn't parse request:\n\"%s\"", request)
	}

	args = split[1]

	routeName = split[0]

	// Find the route's sequence constructor.
	rtFn, ok := routes[routeName]
	if !ok {
		return Route{}, fmt.Errorf("route name \"%s\" not recognised", routeName)
	}
	rt := rtFn()

	// Unmarshall the toml arguments, if they were given.
	if args != "" {
		_, err := toml.Decode(args, rt.Params)
		if err != nil {
			return Route{}, fmt.Errorf("couldn't unmarshall arguments (%s)", err.Error())
		}
	}

	return rt, nil
}

func (rt *Route) Sequence() sequence.RunAller {
	return rt.NewSequence(rt.Params)
}

type Params interface {
	IsParams()
}

