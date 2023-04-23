package db

import(
    "fmt"
    "errors"
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
