package graph_model

import (
    "fmt"
)


func (csvg * CSV_Graph) String_from_all_nodes(read_nodes map[uint32]bool) (string, int) {
    nodes_string := ""
    for _, node := range csvg.Nodes{
        if visited, exists := read_nodes[node.ID]; !exists || !visited{
            nodes_string += node.String() + "\n"
            read_nodes[node.ID] = true
        }
        continue
    }
    return nodes_string,  len(read_nodes)
}


func (csvg * CSV_Graph) String_from_all_edges(read_nodes map[uint32]bool) (string, int){
    edges_count := 0
    edges_string := ""
    for _, node := range csvg.Nodes {
        if visited, exists := read_nodes[node.ID]; !visited || !exists {
            for _, neighbour_id := range node.NeighboursID {
                read_nodes[node.ID] = true
                edges_string += fmt.Sprintf("%v\t%v\n", node.ID, neighbour_id)
                edges_count++
            }
        }
    }
    return edges_string, edges_count
}
