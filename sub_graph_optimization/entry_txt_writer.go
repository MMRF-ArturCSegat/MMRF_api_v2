package sub_graph_optimization

import (
	"fmt"
    gm "github.com/UFSM-Routelib/routelib_api/graph_model"
    "os"
)

// type to deal with the return value of gm.CSV_Graph.SpreadRadius()
// not ideal representation of a SubGraph of a gm.CSV_Graph
// maybe refactor gm.CSV_Graph.SpreadRadius to return a CSV_Graph would be a great idea
// maybe a tree structure would also work well for this purposes
type SubGraph []gm.GraphPath


func (sg SubGraph) all_nodes_string() string {
    read_nodes := make(map[int64]bool)
    nodes_string := ""
    for _, path := range sg{
        for _, node := range path.Nodes{
            if visited, exists := read_nodes[node.ID]; !exists || !visited{
                nodes_string += node.String() + "\n"
                read_nodes[node.ID] = true
            }
            continue
        }
    }
    return nodes_string
}


func (sg SubGraph) all_edges_string() string {
    read_nodes := make(map[int64]bool)
    edges_string := ""
    for _, path := range sg {
        for _, node := range path.Nodes {
            if visited, exists := read_nodes[node.ID]; !visited || !exists {
                for _, neighbour_id := range node.NeighboursID {
                    read_nodes[node.ID] = true
                    edges_string += fmt.Sprintf("%v\t%v\n", node.ID, neighbour_id)
                }
            }
        }
    }
    return edges_string
}

func (sg SubGraph) Build_txt_file() (*os.File,   error) {
    file_content := ""
    file_content += "Nodes\n"
    file_content += sg.all_nodes_string()
    file_content += "Edges\n"
    file_content += sg.all_edges_string()
    file, err := os.Create("sub_graph.txt")
    if err != nil {
        return file, err
    }
    _, err = file.WriteString(file_content)
    if err != nil {
        return file, err
    }
    return file, nil
}
