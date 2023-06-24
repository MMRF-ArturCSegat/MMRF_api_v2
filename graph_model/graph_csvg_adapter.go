package graph_model

import "github.com/UFSM-Routelib/routelib_api/util"

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


// stuplidly inefficient to draw on the map
func (csvg * CSV_Graph) Csvg_to_slice_of_coord_paths() [][]util.Coord{
    all_nodes := csvg.AllNodes()
    paths := make([][]util.Coord, 0, len(all_nodes))
    for _, node := range all_nodes{
        for _, neigh_id := range node.NeighboursID{
            neigh, _ := csvg.FindNode(neigh_id)
            paths = append(paths, []util.Coord{node.GetCoord(), neigh.GetCoord()})
        }
    }
    return paths
}
