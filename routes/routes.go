package routes

import (
	"github.com/nedp/remotecmds/router"

	"github.com/nedp/remotecmds/say"
	"github.com/nedp/remotecmds/utc"
)

func NoParams() router.Params {
	return *new(router.Params)
}

func AddRoutesTo(r router.Interface) {
	r.AddRoute("say", say.NewSequence, say.NewParams)
	r.AddRoute("utc", utc.NewSequence, utc.NewParams)
}
