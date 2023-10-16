package graph_model

import (
	"errors"
	"math"
	"slices"
	"sync"
	"github.com/UFSM-Routelib/routelib_api/util"
)

func (csvg * CSV_Graph) LimitedBranchigFrom(start *GraphNode, limit float32, path GraphPath, paths []GraphPath) []GraphPath {
    path.Append(start)

	r := make(chan []GraphPath, len(start.NeighboursID))      // the channel and the waitgroup grantee the children will send their paths aproproatly
	wg := new(sync.WaitGroup)

    valid_neighbours := make([]*GraphNode, 0)

	for _, node_id := range start.NeighboursID{
        if util.In(node_id, path.IdSlice()){ // Validates no backtracking
            continue
        }

        node, err := csvg.FindNode(node_id)
        if err != nil {continue}

        node_coord := node.GetCoord()

        if (path.Cost + node_coord.DistanceToInMeters(start.GetCoord())) > limit{// Validates the node is in the disected square             
            continue
        }

        valid_neighbours = append(valid_neighbours, node)
	}
    
    if len(valid_neighbours) == 1{                                                  
        r <- csvg.LimitedBranchigFrom(valid_neighbours[0], limit, path.Copy(), make([]GraphPath, 0))
    }
    if len(valid_neighbours) > 1 {                                                  // Theses two if's generate the selective threading
                                                                                    // We will only crate a new thread if the node
        for _, valid_node := range valid_neighbours{                                // Has more than two valide neighbours, otherwise
			go func(node *GraphNode){                                               // Just keep the same thread computing              
                defer wg.Done()

                subPaths := csvg.LimitedBranchigFrom(node, limit, path.Copy(), make([]GraphPath, 0)) 

                r <- subPaths 
            }(valid_node)

			wg.Add(1)
        }
    }
	
	wg.Wait()
	close(r)

	paths = append(paths, path)                     // This section add the start path to the paths
	for subPaths := range r{                        // Then ummount the channel to receive the paths from the children ann add them to the paths
		for i := 0; i< len(subPaths); i++{
			paths = append(paths, subPaths[i])
		}
	}

	return paths
}

func (csvg * CSV_Graph) ClosestNode(co util.Coord) (*GraphNode, float32){
    type result struct{     // type for distance calulation
        node    *GraphNode
        dist    float32
        mut     *sync.Mutex
    }

    // worker function to calculate a section of the array
    computer := func(from, to int, arr []*GraphNode, best_global *result, wg * sync.WaitGroup){
        defer wg.Done()
        var best_node *GraphNode
        best_node = nil
        best_dist := float32(math.MaxFloat32)

        for _, node := range arr[from:to] {
            if d := node.GetCoord().DistanceToInMeters(co); d < best_dist || best_dist == 0 {
                best_dist = d
                best_node = node
            }
        }
        
        best_global.mut.Lock()
        if best_dist < (*best_global).dist {
            (*best_global).node = best_node
            (*best_global).dist = best_dist
        }
        best_global.mut.Unlock()
    }

    nodes := csvg.AllNodes()
    wg := new(sync.WaitGroup)
    best_result := result{node: nodes[0], dist: nodes[0].GetCoord().DistanceToInMeters(co), mut: new(sync.Mutex)}

    // weird foor loop call the worker for each 10% of the nodes slice
	step_c := len(nodes)/10
	step_0 := 0
	step := step_c
    for i := 0; i<10; i++{
        go computer(step_0, step, nodes, &best_result, wg)
        step_0 += step_c        
        step += step_c
        wg.Add(1)
    }
    wg.Wait()

    return best_result.node, best_result.dist   // unmounting of best result for no type conflicts
}

// this function is implemented syncronosly as to avoid race condition
// with the comp func and array reading
// an example would be in instance.go/Instance.OltNecessaryConnections
func (csvg * CSV_Graph) ClosestNodeFunc(co util.Coord, comp func(node *GraphNode, reference util.Coord, dist float32)bool ) (*GraphNode, float32, error){
    type result struct{     // type for distance calulation
        node    *GraphNode
        dist    float32
    }
    
    nodes := csvg.AllNodes()
    // this is different from the default ClosestNode so that the first node from AllNodes() is not ignored
    best_result := result{node: nil, dist: math.MaxFloat32}

    for _, node := range nodes{
        if d := node.GetCoord().DistanceToInMeters(co); comp(node, co, best_result.dist) || best_result.dist == 0 {
            best_result.dist = d
            best_result.node = node
        }

    }

    // no node in AllNodes() was abble to passes the filter set in the comp func
    if best_result.node == nil || best_result.dist == math.MaxFloat32 {
        return nil, 0.0, errors.New("No node passed your defined filter")
    }

    return best_result.node, best_result.dist, nil  // unmounting of best result for no type conflicts
}

func (csvg * CSV_Graph) order_neighbours_by_prox_to(ref util.Coord, neigh_id_slice []uint32) []*GraphNode{
    ns := make([]*GraphNode, 0, len(neigh_id_slice))
    for _, id := range neigh_id_slice {
        if n, err := csvg.FindNode(id); err == nil{
            ns = append(ns, n)
        }
    }
        
    slices.SortFunc(ns, func(a, b *GraphNode)int{
        return int(a.GetCoord().DistanceToInMeters(ref) - b.GetCoord().DistanceToInMeters(ref))
    })
    
    return ns
}

func (csvg * CSV_Graph)walk(start, end *GraphNode, path *GraphPath, seen *map[uint32]bool) bool{
    path.Append(start)

    if start.ID == end.ID {
        return true
    }

    if _, ok := (*seen)[start.ID]; ok {
        path.Pop()
        return false
    }
    (*seen)[start.ID] = true

    for _, node := range csvg.order_neighbours_by_prox_to(end.GetCoord(), start.NeighboursID) {
        c := path.Copy()
        if csvg.walk(node, end, &c, seen) {
            *path = c
            return true
        }
    }

    path.Pop()
    return false
}

func (csvg * CSV_Graph) DepthFirstSearch(start, end *GraphNode) GraphPath {
    path := GraphPath{}
    seen := make(map[uint32]bool, len(csvg.Nodes)/2)
    
    csvg.walk(start, end, &path, &seen)
    
    if len(path.Nodes) == 1 && start.ID != end.ID {
        return GraphPath{}
    }
    return path 
}

func (csvg * CSV_Graph) AllVisitableNodesFrom(start *GraphNode, fill []*GraphNode) []*GraphNode {
    fill = append(fill, start)
    for _, node_id := range start.NeighboursID {
        node, err := csvg.FindNode(node_id)
        if err != nil {continue}
        if !util.In(node, fill) {
            fill = csvg.AllVisitableNodesFrom(node, fill)
        }
    }
    return fill
}

func (csvg * CSV_Graph) AllNodes() []*GraphNode{

    nodes := make([]*GraphNode, 0, len(csvg.Nodes))
    for _, value := range csvg.Nodes{
        nodes = append(nodes, value)
    }
    return nodes
}
