package main

import (
	"github.com/UFSM-Routelib/routelib_api/routes"
)

func main(){
	r := routes.SetupRouter()
	r.Run(":1337")
}
