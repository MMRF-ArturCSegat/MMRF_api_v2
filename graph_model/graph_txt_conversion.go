package graph_model

import (
    "os"
    "fmt"
)


func (csvg * CSV_Graph) all_nodes_string() string {
    read_nodes := make(map[int64]bool)
    nodes_string := ""
    for _, node := range csvg.Nodes{
        if visited, exists := read_nodes[node.ID]; !exists || !visited{
            nodes_string += node.String() + "\n"
            read_nodes[node.ID] = true
        }
        continue
    }
    return nodes_string
}


func (csvg * CSV_Graph) all_edges_string() string {
    read_nodes := make(map[int64]bool)
    edges_string := ""
    for _, node := range csvg.Nodes {
        if visited, exists := read_nodes[node.ID]; !visited || !exists {
            for _, neighbour_id := range node.NeighboursID {
                read_nodes[node.ID] = true
                edges_string += fmt.Sprintf("%v\t%v\n", node.ID, neighbour_id)
            }
        }
    }
    return edges_string
}

func (csvg * CSV_Graph) Build_txt_file() (*os.File,   error) {
    file_content := ""
    file_content += "Nodes\n"
    file_content += csvg.all_nodes_string()
    file_content += "Edges\n"
    file_content += csvg.all_edges_string()
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
