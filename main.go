package main

import (
	"github.com/UFSM-Routelib/routelib_api/routes"
	"github.com/UFSM-Routelib/routelib_api/sessions"
)

func main(){
	r := routes.SetupRouter()
    go sessions.CleanExpiredSessions()
	r.Run(":1337")
}
