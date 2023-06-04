package graph_model

import (
    "github.com/UFSM-Routelib/routelib_api/util"
    "sync"
)

func (csvg * CSV_Graph) SpreadRadius(start *GraphNode, limit float32, path GraphPath, paths []GraphPath, square util.Square) []GraphPath {
    path.Append(start)

	r := make(chan []GraphPath, len(start.NeighboursID))      // the channel and the waitgroup grantee the children will send their paths aproproatly
	wg := new(sync.WaitGroup)

    valid_neighbours := make([]*GraphNode, 0)

	for _, node_id := range start.NeighboursID{
        if util.In(node_id, path.IdSlice()){ // Validates no backtracking
            continue
        }

        node, _ := csvg.FindNode(node_id)

        node_coord := node.GetCoord()

        if !node_coord.IsInSquare(square){   // Validates the node is in the disected square             
            continue
        }

        if (path.Cost + node_coord.DistanceToInMeters(start.GetCoord())) > limit{// Validates the node is in the disected square             
            continue
        }

        valid_neighbours = append(valid_neighbours, node)
	}
    
    if len(valid_neighbours) == 1{                                                  
        r <- csvg.SpreadRadius(valid_neighbours[0], limit, path.Copy(), make([]GraphPath, 0), square)
    }
    if len(valid_neighbours) > 1 {                                                  // Theses two if's generate the selective threading
                                                                                    // We will only crate a new thread if the node
        for _, valid_node := range valid_neighbours{                                // Has more than two valide neighbours, otherwise
			go func(node *GraphNode){                                                    // Just keep the same thread computing              
                defer wg.Done()

                subPaths := csvg.SpreadRadius(node, limit, path.Copy(), make([]GraphPath, 0), square) 

                r <- subPaths 
            }(valid_node)

			wg.Add(1)
        }
    }
	
	wg.Wait()
	close(r)

	paths = append(paths, path)                     // This section add the current path to the paths
	for subPaths := range r{                        // Then ummount the channel to receive the paths from the children ann add them to the paths
		for i := 0; i< len(subPaths); i++{
			paths = append(paths, subPaths[i])
		}
	}

	return paths
}

func (csvg * CSV_Graph) ClosestNode(co util.Coord) (*GraphNode, float32){
    type result struct{     // type for distance calulation
        node *GraphNode
        dist float32
    }

    // worker function to calculate a section of the array
    computer := func(from, to int, arr []*GraphNode, results chan result, wg * sync.WaitGroup){
        defer wg.Done()
        best_node := arr[0]
        best_dist := best_node.GetCoord().DistanceToInMeters(co)
        for from < to {
            if c := arr[from].GetCoord(); c.DistanceToInMeters(co) < best_dist || best_dist == 0{
                best_dist = c.DistanceToInMeters(co)
                best_node = arr[from]
            }
            from += 1
        }
        result := result{node: best_node, dist: best_dist}
        results <- result
    }

    nodes := csvg.AllNodes()
    results :=  make(chan result, 10)
    wg := new(sync.WaitGroup)

	step_c := len(nodes)/10
	step_0 := 0
	step := step_c
    for i := 0; i<10; i++{
        println("s0", step_0, "sp", step)
        go computer(step_0, step, nodes, results, wg)
        step_0 += step_c        // weird foor loop call the worker for each 10% of the nodes slice
        step += step_c
        wg.Add(1)
    }

    wg.Wait()
    close(results)

    best_result := <- results

    for result := range results {           // collection and calulation
        if result.dist < best_result.dist{
            best_result = result
        }
    }

    return best_result.node, best_result.dist   // unmounting of best result for no type conflicts
}


func (csvg * CSV_Graph) AllNodes() []*GraphNode{
    nodes := make([]*GraphNode, 0, len(csvg.Nodes))
    for _, value := range csvg.Nodes{
        nodes = append(nodes, value)
    }
    return nodes
}
