package graph_model

import (
	"fmt"
	"github.com/UFSM-Routelib/routelib_api/util"
)


type GraphNode struct {
    ID		        uint32		`json:"id"`
	Lat		        float64		`json:"lat"`
	Lng		        float64		`json:"lng"`
    NeighboursID    []uint32     `json:"neighbours"`
}


func (n * GraphNode) GetCoord() util.Coord{
    return util.NewCoord(n.Lat, n.Lng)          // I used this because i was lazy to refactor the whole api
}                                               // Instead Nodes will have those two atributtes and when you need 
                                                //to useit as a Coord type just call the builder

func (n GraphNode) Debug () string {
    return fmt.Sprintf(" node[%v, %v, %v, %v] ", n.ID, n.Lat, n.Lng, n.NeighboursID)
}


func (n * GraphNode) String() string {
    return fmt.Sprintf("%v\t%v\t%v", n.ID, n.Lat, n.Lng)
}









