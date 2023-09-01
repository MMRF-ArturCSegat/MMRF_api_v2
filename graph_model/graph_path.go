package graph_model

import (
	"fmt"
)

type GraphPath struct{
    Nodes           []*GraphNode
    Cost            float32
}

func (p * GraphPath) NodeIn(node_id uint32) bool {
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
    if len(p.Nodes) >= 1{                       // the first node, should have no cost to traverse
        last_node := p.Nodes[len(p.Nodes) - 1]
        coord := last_node.GetCoord()
        dist := coord.DistanceToInMeters(n.GetCoord())
        p.Cost += dist
    }
    p.Nodes = append(p.Nodes, n)
}

func (p * GraphPath) Pop() {                   
    if len(p.Nodes) == 0 {
        return 
    }

    increment := float32(0.0)
    if len(p.Nodes) > 1 {
        increment = p.Nodes[len(p.Nodes) - 1].GetCoord().DistanceToInMeters(p.Nodes[len(p.Nodes) - 2].GetCoord())
    }
    p.Nodes = p.Nodes[len(p.Nodes) - 1:]
    p.Cost -= increment
}

// useful so children can modify their parent's path, without chaningin common underlying array
func (p * GraphPath) Copy() GraphPath {                                 
    b := GraphPath{Nodes: make([]*GraphNode, len(p.Nodes)), Cost: 0}   
    b.Cost = p.Cost                                                   

    for i, n := range p.Nodes {
        if n == nil {
            continue
        }
        v := *n
        b.Nodes[i] = &v
    }
    
    return b
}

func (p * GraphPath) Print(){
    fmt.Printf("nodes: %v | Cost: %f \n", p.IdSlice(), p.Cost)
}
