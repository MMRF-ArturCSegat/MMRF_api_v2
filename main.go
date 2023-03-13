package main

import (
	"gat/test"
	"gat/api_funcs"
)

func init(){
	db2.ConnectDatabase2()
}

func main(){
	r := funcs.SetupRouter()
	r.Run(":3000")
}
