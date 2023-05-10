package db

import (
	"errors"
	"fmt"

	"github.com/MMRF-ArturCSegat/MMRF_api_v2/util"
)

func AddNode(n *Node) (*Node, error){
	fmt.Println("Adding node " ,n)
	//	This function assumes that the Node does not exist in the database as the checker to see if it exists happens in the API route	
	res := db.Create(&n)

	if res.RowsAffected == 0{
		return nil, errors.New("Couldn't add node to database")
	}
	
	return n, nil
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
    err := ConnectNodes(node1, node2)

    if err != nil{
        return nil, errors.New("Failed to connect both nodes")
    }

	return []*Node{node1, node2}, nil
}

func ConnectNodes(n1, n2 *Node) error {
	fmt.Println("Connecting nodes")
    
    if !util.In(n2.ID, n1.Neighbours){
        val1 := append(n1.Neighbours, n2.ID)
        fmt.Println("val1",val1)
        if err := db.Model(n1).UpdateColumn("neighbours", val1).Error; err != nil {
            return errors.New("Failed to connect on node1")
        }
    }

    if !util.In(n1.ID, n2.Neighbours){
        val2 := append(n2.Neighbours, n1.ID)
        fmt.Println("val2",val2)
        if err := db.Model(n2).UpdateColumn("neighbours", val2).Error; err != nil {
            return errors.New("Failed to connect on node2")
        }
    }

	return nil
}
