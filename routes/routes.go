package routes

import (
	"bitbucket.org/nedp/remotecmds/router"

	"bitbucket.org/nedp/remotecmds/say"
	"bitbucket.org/nedp/remotecmds/utc"
)

func NoParams() router.Params {
	return *new(router.Params)
}

func AddRoutesTo(r router.Interface) {
	r.AddRoute("say", say.NewSequence, say.NewParams)
	r.AddRoute("utc", utc.NewSequence, NoParams)
}
