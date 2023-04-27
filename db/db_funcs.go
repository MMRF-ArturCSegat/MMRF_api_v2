package db

import(
    "fmt"
    "errors"
    "gat/utilities"
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

func AllNodes() ([]*Node, error){
	var nodes []*Node

	res := db.Find(&nodes)

	if res.Error != nil{
		return nil, errors.New("Failed to find all nodes")
	}
	return nodes, nil
}

func IdSliceFromNodeSlice(node_slice []*Node) []int64 {
	id_slice := make([]int64, len(node_slice))

	if len(node_slice) == 0{
		return id_slice
	}

	for i, e := range node_slice{
		id_slice[i] = e.ID	
	}
	return id_slice
}

func SpreadRadius(start *Node, limit, cost int, path []*Node, paths [][]*Node) [][]*Node {
	path = append(path, start)
	cost += 1

	dissectPath := func(node *Node, limit, cost int, path []*Node, results chan [][]*Node, wg * sync.WaitGroup) { // function will be used to fill a channe with paths
		defer wg.Done()
	
		subPaths := SpreadRadius(node, limit, cost, path, make([][]*Node, 0)) // calculating subpaths by dividing the graph in subgraphs

		results <- subPaths // sends the resulting paths into the chanel so they can be read later
	}

	r := make(chan [][]*Node, len(start.Neighbours))
	wg := new(sync.WaitGroup)

	routine_counter := 0

	for _, node_id := range start.Neighbours{
		if !util.In(node_id, IdSliceFromNodeSlice(path)){
			if cost > limit{ // if here adding the next node would not pass the limit
				continue
			}
			node, _ := FindNode(node_id)

			p := make([]*Node, len(path))
			copy(p, path)
			go dissectPath(node, limit, cost, p, r, wg)
			wg.Add(1)
			routine_counter += 1
		}
	}
	
	wg.Wait()
	close(r)

	paths = append(paths, path)
	for subPaths := range r{
		for i := 0; i< len(subPaths); i++{
			paths = append(paths, subPaths[i])
		}
	}

	return paths
}
