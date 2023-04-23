package main

import (
	"gat/db"
	"gat/routes"
)

func init(){
	db.ConnectDatabase2()
}

func main(){
	r := routes.SetupRouter()
	r.Run(":1337")
}
