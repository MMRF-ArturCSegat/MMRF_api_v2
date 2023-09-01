package graph_model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newNode(id uint32, lat, lng float64) *GraphNode {
    return &GraphNode{
        ID: id,
        Lat: lat,
        Lng: lng,
        NeighboursID: make([]uint32, 0),
    }
}

func setup() *CSV_Graph {
    c := &CSV_Graph{Nodes: make(map[uint32]*GraphNode)}

    c.AddEdge(newNode(1, 1, 1), newNode(2, 2, 2))
    c.AddEdge(newNode(1, 1, 1), newNode(3, 3, 3))
    c.AddEdge(newNode(3, 3, 3), newNode(4, 4, 4))
    c.AddEdge(newNode(2, 2, 2), newNode(4, 4, 4))
    c.AddEdge(newNode(2, 2, 2), newNode(5, 5, 5))

    return c
}
    
func TestAllVisitableNodes(t * testing.T) {
    c := setup()   
    
    n1, _ := c.FindNode(1)
    n2, _ := c.FindNode(2)
    n3, _ := c.FindNode(3)
    n4, _ := c.FindNode(4)
    n5, _ := c.FindNode(5)
    nodes := make([]*GraphNode, 0)
        

    if av := c.AllVisitableNodesFrom(n1, nodes);!assert.Equal(t, []*GraphNode{n1, n2, n4, n3, n5}, av) {
        fmt.Println(av)
        fmt.Println([]*GraphNode{n1, n2, n4, n3, n5})
        t.Errorf("exepcted %v, got %v", []*GraphNode{n1, n2, n4, n3, n5}, av)
        return
    }
    println("Passed AllVisitableNodes")
}


func TestDFS(t * testing.T) {
    c := setup()

    n1, _ := c.FindNode(1)
    n2, _ := c.FindNode(2)
    n4, _ := c.FindNode(4)
    
    dfs := c.DepthFirstSearch(n1, n4)
    distance := n1.GetCoord().DistanceToInMeters(n2.GetCoord()) + n2.GetCoord().DistanceToInMeters(n4.GetCoord())
    
    expect := &GraphPath{Nodes: []*GraphNode{n1, n2, n4}, Cost: distance}

    if !assert.Equal(t, expect, dfs) {
        t.Errorf("1 -> 4: expected %v, got %v", expect, dfs)
        return
    }

    c.AddNode(newNode(6, 6, 6))
    n6, _ := c.FindNode(6)

    expect = nil
    dfs = c.DepthFirstSearch(n1, n6)

    if !assert.Equal(t, expect, dfs) {
        t.Errorf("1 -> 6: expected %v, got %v", expect, dfs)
        return
    }
    println("Passed DFS")
}   
