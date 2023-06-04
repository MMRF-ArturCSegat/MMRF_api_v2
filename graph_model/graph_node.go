
package graph_model

import (
    "github.com/UFSM-Routelib/routelib_api/util"
)


type GraphNode struct {
    ID		        int64		`json:"id"`
	Lat		        float64		`json:"lat"`
	Lng		        float64		`json:"lng"`
    NeighboursID    []int64     `json:"neighbours"`
}


func (n * GraphNode) GetCoord() util.Coord{
    return util.NewCoord(n.Lat, n.Lng)          // I used this because i was lazy to refactor the whole api
                                                // for nodes to have an Coord atribute of Coord type
}                                               // Instead Nodes will have those two atributtes and when you need 
                                                //to useit as a Coord type just call the builder









