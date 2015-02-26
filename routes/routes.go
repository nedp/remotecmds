package routes

import (
	"bitbucket.org/nedp/remotecmds/router"
	"bitbucket.org/nedp/remotecmds/say"
)

func Routes(name string) (route router.Route, ok bool) {
	ok = true
	switch name {
	case "say":
		route = router.Route{
			"say",
			say.New,
			new(say.Params),
		}
	default:
		ok = false
	}
	return
}
