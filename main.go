package main

import (
	"gat/database"
	"gat/api_funcs"
)

func init(){
	db.ConnectDatabase()
}

func main(){
	r := funcs.SetupRouter()
	r.Run(":3000")
}