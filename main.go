package main

import (
	"time"
	foc"github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
	"github.com/UFSM-Routelib/routelib_api/routes"
	"github.com/UFSM-Routelib/routelib_api/sessions"
)


func session_cleaner(){
    for {
        time.Sleep(15 * time.Minute)
        sessions.CleanExpiredSessions()
    }
}


func init(){
    foc.ConnectFBCDatabase()
}


func main(){
	r := routes.SetupRouter()
    go session_cleaner()
	r.Run(":1337")
}
