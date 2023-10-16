package graph_model

import (
	"fmt"
	"testing"
	"github.com/UFSM-Routelib/routelib_api/util"
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
    c := &CSV_Graph{
        Nodes: make(map[uint32]*GraphNode),
        Limiter: util.DefaultMaxSquare(),
        Olt: util.NewCoord(10, 10),
    }

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
    fmt.Println("Passed AllVisitableNodes")
}


// func TestDFS(t * testing.T) {
//     c := setup()
//
//     n1, _ := c.FindNode(1)
//     n2, _ := c.FindNode(2)
//     n4, _ := c.FindNode(4)
//     
//     dfs := c.DepthFirstSearch(n1, n4)
//     distance := n1.GetCoord().DistanceToInMeters(n2.GetCoord()) + n2.GetCoord().DistanceToInMeters(n4.GetCoord())
//     
//     expect := &GraphPath{Nodes: []*GraphNode{n1, n2, n4}, Cost: distance}
//
//     if !assert.Equal(t, expect, dfs) {
//         t.Errorf("1 -> 4: expected %v, got %v", expect, dfs)
//         return
//     }
//
//     c.AddNode(newNode(6, 6, 6))
//     n6, _ := c.FindNode(6)
//
//     expect = nil
//     dfs = c.DepthFirstSearch(n1, n6)
//
//     if !assert.Equal(t, expect, dfs) {
//         t.Errorf("1 -> 6: expected %v, got %v", expect, dfs)
//         return
//     }
//     fmt.Println("Passed DFS")
// }   

func TestCloestNode(t * testing.T) {
    c := setup()

    for i := 1; i <= 10; i++ {    
        c.AddNode(newNode(uint32(i), float64(i), float64(i)))
    }

    n1, _ := c.FindNode(1)
    
    expected_node := n1
    expected_dist := n1.GetCoord().DistanceToInMeters(util.NewCoord(1.2, 1.2))

    node, dist := c.ClosestNode(util.NewCoord(1.2, 1.2))

    if node != expected_node || dist != expected_dist {
        t.Errorf("Bad distance or node from coordinate {1.2, 1.2}: expected node %v with distance %v, got node %v and distance %v", 
            expected_node, expected_dist, node, dist)
            return
    }
    fmt.Println("Passed ClosestNode")
}


func TestCloestNodeFunc(t * testing.T) {
    c := setup()

    for i := 1; i <= 10; i++ {    
        c.AddNode(newNode(uint32(i), float64(i), float64(i)))
    }

    n1, _ := c.FindNode(1)
    
    expected_node := n1
    expected_dist := n1.GetCoord().DistanceToInMeters(util.NewCoord(1.2, 1.2))
    
    filter1 := func (node *GraphNode, reference util.Coord, dist float32)bool {
        if node.GetCoord().DistanceToInMeters(reference) < dist {
            return true
        }
        return false
    }

    node, dist, err := c.ClosestNodeFunc(util.NewCoord(1.2, 1.2), filter1)

    if err != nil {
        t.Errorf("error happend with coordinate {1.2, 1,2} filter 1: %v", err.Error())
        return
    }

    if node != expected_node || dist != expected_dist {
        t.Errorf("Bad distance or node from coordinate {1.2, 1.2} filter 1: expected node %v with distance %v, got node %v and distance %v", 
            expected_node, expected_dist, node, dist)
            return
    }


    filter2 := func (node *GraphNode, reference util.Coord, dist float32)bool {
        if node.GetCoord().DistanceToInMeters(reference) < dist && node.ID > 10{
            return true
        }
        return false
    }

    node, dist, err = c.ClosestNodeFunc(util.NewCoord(1.2, 1.2), filter2)
        
    expected_node = nil
    expected_dist = 0.0

    if err == nil {
        t.Errorf("The function should error like :'No node passed your defined filter' but it showed no error; Problem on filter 2")
        return
    }

    if node != expected_node || dist != expected_dist {
        t.Errorf("Bad distance or node from coordinate {1.2, 1.2} filter 2: expected node %v with distance %v, got node %v and distance %v", 
            expected_node, expected_dist, node, dist)
            return
    }
    fmt.Println("Passed ClosestNodeFunc")
}


func TestOltNecessaryConnections(t * testing.T) {
    g := setup()
    
    t1 := newNode(12, 12, 12)
    t2 := newNode(13, 13, 13)
    t3 := newNode(14, 14, 14)

    g.AddEdge(t1, t2)
    g.AddEdge(t2, t3)

    n_conns := g.OltNecessaryConnections()
    
    if len(n_conns) != 2 {
        t.Errorf("Bad connections len: len should be 2 is %v\n", len(n_conns))
        return
    }

    fill_1 := make([]*GraphNode, 0)
    fill_2 := make([]*GraphNode, 0)
    
    n1, _ := g.FindNode(n_conns[0])
    n2, _ := g.FindNode(n_conns[1])

    fill_1 = g.AllVisitableNodesFrom(n1, fill_1)
    fill_2 = g.AllVisitableNodesFrom(n2, fill_2)

    if len(fill_1) + len(fill_2) != len(g.AllNodes()) {
        t.Errorf("Bad connections: shoul dbe able to reach all nodes(total %v), can only reach %v\n", len(g.AllNodes()), len(fill_1) + len(fill_2))
        fmt.Printf("fill_1: %v\n", fill_1)
        fmt.Printf("fill_2: %v\n", fill_2)
        return
    }
    
    fmt.Println("Passed OltNecessaryConnections")
}

