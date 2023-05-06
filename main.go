package main

import (
	"github.com/MMRF-ArturCSegat/MMRF_api_v2/db"
	"github.com/MMRF-ArturCSegat/MMRF_api_v2/routes"
)

func init(){
	db.ConnectDatabase2()
}

func main(){
	r := routes.SetupRouter()
	r.Run(":1337")
}
