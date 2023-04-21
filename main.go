package main

import (
	"gat/database"
	"gat/api_funcs"
)

func init(){
	db.ConnectDatabase2()
}

func main(){
	r := funcs.SetupRouter()
	r.Run(":3000")
}
