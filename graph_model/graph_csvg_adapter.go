package graph_model

func Slice_of_paths_to_csvg(paths []GraphPath) *CSV_Graph {
    csv_graph := CSV_Graph{Nodes: make(map[uint32]*GraphNode)}
    
    for _, path := range paths{
        Inner: 
        for index, node := range path.Nodes{
            if index == 0 {
                continue Inner
            }
            csv_graph.AddEdge(node, path.Nodes[index - 1])
            if index == len(path.Nodes) - 1{
                continue Inner
            }
            csv_graph.AddEdge(node, path.Nodes[index + 1])
        }
    }
    return &csv_graph   
}

// TODO: implement this; no real use yet

// func Csvg_to_slice_of_paths(CSV_Graph) []*GraphPath{
// }
