package db

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"fmt"
)

var err error
var db *gorm.DB

type Node struct {
    ID      int       `json:"id" gorm:"primaryKey"`
    Lat     float64   `json:"lat"`
    Lng     float64   `json:"lng"`
    Edges   []*Edge   `json:"edges" gorm:"foreignKey:EID"	`
}

type Edge struct {
	EID     	int		`json:"id" gorm:"primaryKey"`
	Node1ID		int     	`json:"node1_id"`
	Node2ID		int		`json:"node2_id"`
	Cost    	int     	`json:"cost"`
}

func contains(list []*Edge, x int) bool {
	for _, item := range list {
		if item.EID == x {
			return true
		}
	}
	return false
}

func ConnectDatabase() {
    // Replace the connection details below with your own MySQL database configuration
	dsn := "arturcs:123123123@tcp(127.0.0.1:3306)/gat_db?charset=utf8mb4&parseTime=True&loc=Local"
    database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }

    err = database.AutoMigrate(&Node{}, &Edge{})
    if err != nil {
        fmt.Println("Failed to auto migrate database schema")
        return
    }


    db = database
}

func FindNode(id int) (*Node, error){
	var node Node
	res := db.First(&node,"id = ?", id)
	if res.RowsAffected == 0 || res.Error != nil {
		fmt.Println("Node is nil")
		fmt.Println(node)
		return nil, errors.New("Not in database")
	} 

	return &node, nil

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

	// eid  := strconv.Itoa(n1.ID) +  "." + strconv.Itoa(n2.ID) 		// form the Edge ID in a cool way
	// rev_eid := strconv.Itoa(n2.ID) + "." + strconv.Itoa(n1.ID) 		// reverse so no double edging

	// var ed Edge

	// r := db.First(&ed, "e_id = ?", eid)			// The database is queryed twice so that if A->B exists You cant add
	// r2 := db.First(&ed, "e_id = ?", rev_eid)		


func AddEdge(n1, n2 *Node, cost int) (*Edge, error){ 

	fmt.Println("one-", n1, "two-",n2)

	//	The function assumes both n1, and n2 exist in the database
	//	The check to see if the exist or not and handling of the unhappy path are done sin the API route
		
	e := &Edge{
		Node1ID:	n1.ID,
		Node2ID: 	n2.ID,
		Cost: 		cost,
	}

	if contains(n1.Edges, e.EID) || contains(n2.Edges, e.EID){
		return e, errors.New("Edge already exists")	
	}

	res := db.Create(e)

	if res.RowsAffected == 0{
		return nil, errors.New("Failed to create Edge")
	}

	n1.Edges = append(n1.Edges, e)
	n2.Edges = append(n2.Edges, e)
	db.Save(n1)
	db.Save(n2)
	return e, nil
}

// 																			BIG THOUGHT:
// 			When traversing a graph like this, Start or End are not strongly associated, which means they can be traversed in any order, 
// 			That means that in a Edge with ID 4-19, you would expect to traverse 4->19, and then have to query for a 19-4 edge to traverse the reverse
// 			You can simply check if the Node you are visiting (In this case lets say 19) is the Node1 or Node2 in the Edge and then visit the Node that is not

// TODO IMPLEMENT TRAVERSE


func AllEdges() ([]*Edge, error){
	var edges []*Edge

	res := db.Find(&edges)

	if res.Error != nil{
		return nil, errors.New("Failed to find all Edges")
	}

	return edges, nil
}

func AllNodes() ([]*Node, error){
	var nodes []*Node

	res := db.Find(&nodes)

	if res.Error != nil{
		return nil, errors.New("Failed to find all nodes")
	}

	return nodes, nil
}
