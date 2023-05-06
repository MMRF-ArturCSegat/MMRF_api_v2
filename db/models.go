package db

import (
	pq "github.com/lib/pq"
    "github.com/MMRF-ArturCSegat/MMRF_api_v2/util"
)


type Node struct {
	ID		int64		`json:"id" gorm:"primaryKey"`
	Lat		float64		`json:"lat"`
	Lng		float64		`json:"lng"`
	Neighbours	pq.Int64Array	`json:"neighbours" gorm:"type:integer[]"`  // Dont care to implemnet new table for many to many
}


func (n * Node) GetCoord() util.Coord{
    return util.NewCoord(n.Lat, n.Lng)          // I used this because i was lazy to refactor the whole api for nodes to have an Coord atribute of Coord type
}                                               // Instead Nodes will have those two atributtes and when you need to useit as a Coord type just call the builder









