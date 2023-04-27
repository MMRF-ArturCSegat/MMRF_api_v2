package routes

import(
    "fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"gat/db"
)

func addNode(c *gin.Context){
	var node db.Node

	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in binding": err.Error(),})
		return
	}

	fmt.Println(node)

	res, err := db.AddNode(&node)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error in database:": err.Error(),})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"node added": res})
}


func addEdge(c *gin.Context){

	type Body struct {
		Node1	db.Node 	`jason:"node1"`
		Node2	db.Node 	`jason:"node2"`
	}

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	node1 := body.Node1
	node2 := body.Node2

	fmt.Println("one-", node1, "two-",node2)
	nodes, err := db.AddEdge(&node1, &node2)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes" : nodes,})
}


func connect(c *gin.Context){

	type Body struct {
		Node1	db.Node 	`jason:"node1"`
		Node2	db.Node 	`jason:"node2"`
	}

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	node1 := body.Node1
	node2 := body.Node2

	fmt.Println("one-", node1, "two-",node2)
	err := db.ConnectNodes(&node1, &node2)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes" : []db.Node{node1, node2},})
}
