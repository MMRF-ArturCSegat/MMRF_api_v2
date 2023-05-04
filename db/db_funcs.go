package db

import (
	"errors"
	"fmt"
	"gat/util"
	"sync"
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


func ClosestNode(co util.Coord) *Node{
    nodes, _ := AllNodes()
    var best_node *Node
    var dist float32 
    dist = 0
    
    for _, node := range nodes{
        if c := node.GetCoord(); c.DistanceToInMeters(co) < dist || dist == 0{
            dist = c.DistanceToInMeters(co)
            best_node = node
        }
    }
    return best_node
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

	dissectPath := func(node *Node, limit float32, path GraphPath, results chan []GraphPath, wg * sync.WaitGroup) { // function will be used to fill a channe with paths
		defer wg.Done()
	
		subPaths := SpreadRadius(node, limit, path, make([]GraphPath, 0), square) // calculating subpaths by dividing the graph in subgraphs

		results <- subPaths // sends the resulting paths into the chanel so they can be read later
	}

	r := make(chan []GraphPath, len(start.Neighbours))      // the channel and the waitgroup grantee the children will send their paths aproproatly
	wg := new(sync.WaitGroup)

	for _, node_id := range start.Neighbours{
		if !path.NodeIn(node_id){
			node, _ := FindNode(node_id)
            
            node_coord := node.GetCoord()

            if !node_coord.IsInSquare(square){   // Validates the node is in the disected square             
                continue
            }

            if (path.Cost + node_coord.DistanceToInMeters(start.GetCoord())) > limit{// Validates the node is in the disected square             
                continue
            }

			go dissectPath(node, limit, path.Copy(), r, wg)
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
