package db

import (
	"errors"
	"fmt"
	"sync"

	"github.com/UFSM-Routelib/routelib_api/util"
)

func FindNode(id int64) (*Node, error){
	fmt.Println("FindNode ran for", id)
	var node Node
	res := db.First(&node,"id = ?", id)
	if res.RowsAffected == 0 || res.Error != nil {
		fmt.Println("Node is nil")
		fmt.Println(node)
		return nil, errors.New("Not in database")
	} 
	return &node, nil
}


func ClosestNode(co util.Coord) (*Node, float32){
    type result struct{     // type for distance calulation
        node *Node
        dist float32
    }


    computer := func(from, to int, arr []*Node, results chan result, wg * sync.WaitGroup){      // worker function to calculate a section of the array
        fc := from
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
        fmt.Println("from ", fc, "to ", to , "best was ", result.node, result.dist)
        results <- result
    }

    nodes, _ := AllNodes()
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

        fmt.Println("total best ", best_result.node, best_result.dist)
    return best_result.node, best_result.dist   // unmounting of best result for no type conflicts
}


func AllNodes() ([]*Node, error){
	var nodes []*Node

	res := db.Find(&nodes)

	if res.Error != nil{
		return nil, errors.New("Failed to find all nodes")
	}
	return nodes, nil
}


func SpreadRadius(start *Node, limit float32, path GraphPath, paths []GraphPath, square util.Square) []GraphPath {
    path.Append(start)

	r := make(chan []GraphPath, len(start.Neighbours))      // the channel and the waitgroup grantee the children will send their paths aproproatly
	wg := new(sync.WaitGroup)

    valid_neighbours := make([]*Node, 0)

	for _, node_id := range start.Neighbours{
        if util.In(node_id, path.IdSlice()){ // Validates no backtracking
            continue
        }

        node, _ := FindNode(node_id)

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
        r <- SpreadRadius(valid_neighbours[0], limit, path.Copy(), make([]GraphPath, 0), square)
    }
    if len(valid_neighbours) > 1 {                                                  // Theses two if's generate the selective threading
        routine_counter := 0 // debug var for testing                               // We will only crate a new thread if the node
        for _, valid_node := range valid_neighbours{                                // Has more than two valide neighbours, otherwise
			go func(node *Node){                                                    // Just keep the same thread computing              
                defer wg.Done()

                subPaths := SpreadRadius(node, limit, path.Copy(), make([]GraphPath, 0), square) 

                r <- subPaths 
            }(valid_node)

			wg.Add(1)
            routine_counter += 1
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
