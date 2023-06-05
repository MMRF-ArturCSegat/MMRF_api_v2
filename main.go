package main

import (
	"time"

	"github.com/UFSM-Routelib/routelib_api/routes"
	"github.com/UFSM-Routelib/routelib_api/sessions"
)


func session_cleaner(){
    for {
        time.Sleep(15 * time.Minute)
        sessions.CleanExpiredSessions()
    }
}

func main(){
	r := routes.SetupRouter()
    go session_cleaner()
	r.Run(":1337")
}
