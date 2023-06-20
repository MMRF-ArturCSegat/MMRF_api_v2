package graph_model

import (
	"errors"
	"fmt"
	"github.com/UFSM-Routelib/routelib_api/util"
)

// represents a unweighted undirected graph read from a csv
type CSV_Graph struct{
    Nodes map[uint32]*GraphNode
}


func (csvg * CSV_Graph) AddNode(node *GraphNode){
    csvg.Nodes[node.ID] = node
}


func (csvg * CSV_Graph) AddEdge(n1, n2 *GraphNode){
    node1, err1 := csvg.FindNode(n1.ID)
    node2, err2 := csvg.FindNode(n2.ID)

    if err1 != nil {
        csvg.AddNode(n1)
        node1, _ = csvg.FindNode(n1.ID)
    }
    if err2 != nil{
        csvg.AddNode(n2)
        node2, _ = csvg.FindNode(n2.ID)
    }
    csvg.connectNode(node1, node2)
}


func (csvg * CSV_Graph) connectNode (node1, node2 *GraphNode){
    if !util.In(node2.ID, node1.NeighboursID) && node1.ID != node2.ID{
        node1.NeighboursID = append(node1.NeighboursID, node2.ID)
    }
    if !util.In(node1.ID, node2.NeighboursID) && node2.ID != node1.ID{
        node2.NeighboursID = append(node2.NeighboursID, node1.ID)
    }
}


func (csvg * CSV_Graph) FindNode(id uint32) (*GraphNode, error){
    node, exists := csvg.Nodes[id]

    if !exists{
        return nil, errors.New("node does not exist")
    }
    return node, nil
}


func (csvg CSV_Graph) Print(){
    for node_id, node := range csvg.Nodes{
        fmt.Printf("%v: ", node_id)
        for _, nei := range node.NeighboursID{
            neighbour, _  := csvg.FindNode(nei)
            fmt.Printf(" %v ", neighbour.Debug())
        }
        fmt.Print("\n")
    }
}
