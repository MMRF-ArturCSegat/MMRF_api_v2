
package graph_model

import (
	"fmt"
)

type GraphPath struct{
    Nodes           []*GraphNode
    Cost            float32
}

func (p * GraphPath) NodeIn(node_id uint32) bool {       // self explenaroty why this is usefeul
    for _, np := range p.Nodes{
        if np.ID == node_id{
            return true
        }
    }
    return false
}


func (p * GraphPath) IdSlice() []uint32{                 // mostyle used for debuggin
    return_slice := make([]uint32, 0)
    for _, n := range p.Nodes{
        return_slice = append(return_slice, n.ID)
    }
    return return_slice
}

func (p * GraphPath) Append(n *GraphNode){                   
    if len(p.Nodes) >= 1{ // should not increase the cost less then 2 nodes
        last_node := p.Nodes[len(p.Nodes) - 1]
        coord := last_node.GetCoord()
        dist := coord.DistanceToInMeters(n.GetCoord())
        p.Cost += dist
    }
    p.Nodes = append(p.Nodes, n)
}

func (p * GraphPath) Copy() GraphPath {                                 // important so that when passing a Path down to the children of a node
    b := GraphPath{Nodes: make([]*GraphNode, len(p.Nodes)), Cost: 0}         // the path must be deep copied so that children dont modify each others channels
    b.Cost = p.Cost                                                     // causing weird bugs, the modifying happens as p.Nodes is a slice, wich is a reference to a commom array

    for i, n := range p.Nodes {
        if n == nil {
            continue
        }
        // Create shallow copy of source element
        v := *n
        // Assign address of copy to destination.
        b.Nodes[i] = &v
    }
    
    return b
}

func (p * GraphPath) Print(){
    fmt.Printf("nodes: %v | Cost: %f \n", p.IdSlice(), p.Cost)
}
