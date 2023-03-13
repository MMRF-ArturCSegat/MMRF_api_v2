package db2

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"fmt"
)

var err error
var db *gorm.DB


type Node struct {
	ID		int		`json:"id" gorm:"primaryKey"`
	Lat		float64		`json:"lat"`
	Lng		float64		`json:"lng"`
	Neighbours	[]*Node		`json:":neighbours" gorm:"many2many:node_neighbours;association_jointable_foreignkey:neighbour_id"`
}


func ConnectDatabase2() {
    // Replace the connection details below with your own MySQL database configuration
	dsn := "arturcs:123123123@tcp(127.0.0.1:3306)/gat_db?charset=utf8mb4&parseTime=True&loc=Local"
    database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

func AddShit(n int) (*Node, error){

	node, err := FindNode(n)

	if err != nil {
		return nil, errors.New("Error: Noden not in database")
	}
	node.Neighbours = append(node.Neighbours, &Node{
		Lat: 69.420,
		Lng: 420.69,
	})
	res := db.Model(node).Updates(&node)
	if res.RowsAffected == 0 {
		return nil, errors.New("Error: failed to add shit to node")
	}
	return node, nil
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
func FindNode(id int) (*Node, error){
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
func ConnectNodes(node1, node2 *Node) error {
    fmt.Println("Connecting nodes")

    node1.Neighbours = append(node1.Neighbours, node2)
    if err := db.Model(node1).Association("Neighbours").Append(node2); err != nil {
        return err
    }

    node2.Neighbours = append(node2.Neighbours, node1)
    if err := db.Model(node2).Association("Neighbours").Append(node1); err != nil {
        return err
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
/*
//	When traversing a graph like this, Start or End are not strongly associated, which means they can be traversed in any order, 
// 	That means that in a Edge with ID 4-19, you would expect to traverse 4->19, and then have to query for a 19-4 edge to traverse the reverse
// 	You can simply check if the Node you are visiting (In this case lets say 19) is the Node1 or Node2 in the Edge and then visit the Node that is not

// TODO IMPLEMENT TRAVERSE
*/
