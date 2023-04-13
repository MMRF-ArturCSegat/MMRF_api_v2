package db2

import (
	"errors"
	"fmt"
	pq "github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

var err error
var db *gorm.DB

type Node struct {
	ID		int64		`json:"id" gorm:"primaryKey"`
	Lat		float64		`json:"lat"`
	Lng		float64		`json:"lng"`
	Neighbours	pq.Int64Array	`json:":neighbours" gorm:"type:integer[]"`
}



func ConnectDatabase2() {
    // Replace the connection details below with your own PostgreSQL database configuration
    dsn := "host=localhost user=arturcs password=123123123 dbname=gatdb port=5432 sslmode=disable TimeZone=UTC"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }

    err = database.AutoMigrate(&Node{})
    if err != nil {
        fmt.Println("Failed to auto migrate database schema")
        return
    }

    db = database
}

func AddNode(n *Node) (*Node, error){
	fmt.Println("Adding node " ,n)
	//	This function assumes that the Node does not exist in the database as the checker to see if it exists happens in the API route	
	res := db.Create(&n)

	if res.RowsAffected == 0{
		return nil, errors.New("Couldn't add node to database")
	}
	
	return n, nil
}
func FindNode(id int64) (*Node, error){
	fmt.Println("FindNode ran")
	var node Node
	res := db.First(&node,"id = ?", id)
	if res.RowsAffected == 0 || res.Error != nil {
		fmt.Println("Node is nil")
		fmt.Println(node)
		return nil, errors.New("Not in database")
	} 
	return &node, nil
}



func AddEdge(n1, n2 *Node) ([]*Node ,error){ 	
	fmt.Println("AddEdge ran")

	node1, err1 := FindNode(n1.ID)
	node2, err2 := FindNode(n2.ID)

	if err1 != nil {
		node1, err1  = AddNode(n1)
		if err1 != nil{
			return nil, errors.New("Failed to add first node to db")
		}
	}

	if err2 != nil {
		node2, err2  = AddNode(n2)
		if err2 != nil{
			return nil, errors.New("Failed to add second node to db")
		}
	}
	return []*Node{node1, node2}, nil
}
func ConnectNodes(n1, n2 int64) error {
	fmt.Println("finding nodes")
	node1, err := FindNode(n1)
	if err != nil{
		return err
	}
	fmt.Println("found node 1-", node1)


	node2, err := FindNode(n2)
	if err != nil{
		return err
	}
	fmt.Println("found node 2-", node2)

	fmt.Println("Connecting nodes")

	val1 := append(node1.Neighbours, node2.ID)
	fmt.Println("val1",val1)
	if err := db.Model(node1).UpdateColumn("neighbours", val1).Error; err != nil {
		return errors.New("Failed to connect on node1")
	}

	val2 := append(node2.Neighbours, node1.ID)
	fmt.Println("val2",val2)
	if err := db.Model(node2).UpdateColumn("neighbours", val2).Error; err != nil {
		return errors.New("Failed to connect on node2")
	}

	return nil
}

func AllNodes() ([]*Node, error){
	var nodes []*Node

	res := db.Find(&nodes)

	if res.Error != nil{
		return nil, errors.New("Failed to find all nodes")
	}
	return nodes, nil
}

func In[T comparable](e T, list []T) bool{
	for _, item := range list{
		if item == e{
			return true
		}
	}
	return false
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
	fmt.Println("multi")
	if cost > limit{ // if here adding the next node would not pass the limit
		return paths
	}
	
	path = append(path, start)
	cost += 1
	// current node added to the path, now should continue to its neighbours

	paths = append(paths, path) // it is important we addi to the path instantly so we will have all paths in the radius 

	dissectPath := func(node *Node, limit, cost int, path []*Node, paths [][]*Node, results chan [][]*Node) { // function will be used to fill a channe with paths
	                                                                                                    // from each neighbour of start	
		
		subPaths := SpreadRadius(node, limit, cost, path, make([][]*Node, 0)) // calculating subpaths by dividing the graph in subgraphs

		results <- subPaths // sends the resulting paths into the chanel so they can be read later
	}

	r := make(chan [][]*Node)
	send_counter:= 0 // we need this variable because the if in the following loop makes so we are not sure what number we are sending, it is max len(start.Neighbours)
	for _, node_id := range start.Neighbours{
		if In(node_id, IdSliceFromNodeSlice(path)) == false {
			node, _ := FindNode(node_id) // the error can be ignored because all the neighbours must be real
			fmt.Printf("in was false for %v inside of %v\n", node, IdSliceFromNodeSlice(path))
			send_counter += 1
			go dissectPath(node, limit, cost, path, paths, r)
		}
	}

	for i:=0; i < send_counter; i++ { // here we use the send_counter, because if we just looped len(start.Neighbours) times we might read more than we send
		subPaths := <- r
		for _, p := range subPaths {
			paths = append(paths, p)
		}
	}

	return paths
}

func SpreadRadiusSingle(start *Node, limit, cost int, path []*Node, paths [][]*Node) [][]*Node {
	fmt.Println("single")
	if cost > limit{ // if here adding the next node would not pass the limit
		return paths
	}
	
	path = append(path, start)
	cost += 1
	// current node added to the path, now should continue to its neighbours

	paths = append(paths, path) // it is important we addi to the path instantly so we will have all paths in the radius 

	for _, node_id := range start.Neighbours{
		if In(node_id, IdSliceFromNodeSlice(path)) == false {
			node, _ := FindNode(node_id) // the error can be ignored because all the neighbours must be real
			fmt.Printf("in was false for %v inside of %v\n", node, IdSliceFromNodeSlice(path))
			paths = SpreadRadiusSingle(node, limit, cost, path, paths)
		}
	}
	return paths
}
