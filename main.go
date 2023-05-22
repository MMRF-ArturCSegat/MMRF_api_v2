package main

import (
	"github.com/UFSM-Routelib/routelib_api/db"
	"github.com/UFSM-Routelib/routelib_api/routes"
)

func init(){
	db.ConnectDatabase2()
}

func main(){
	r := routes.SetupRouter()
	r.Run(":1337")
}
