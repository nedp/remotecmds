package routes

import (
	"bitbucket.org/nedp/remotecmds/router"
	"bitbucket.org/nedp/remotecmds/say"
)

func AddRoutesTo(r router.Interface) {
	r.AddRoute("say", say.NewSequence, say.NewParams)
}
