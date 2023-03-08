package funcs

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gat/database"
	"fmt"
)

func home(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message": "Lol funny api",
	})
	return
}

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
	return
}

func addEdge(c *gin.Context){

	type Body struct {
		Nodes 	[]db.Node 	`jason:"nodes"`
		Dist  	int			`jason:"dist"`
	}

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	if len(body.Nodes) != 2{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expecting an array of 2 user objects"})
        return
	}

	node1 := body.Nodes[0] // start and end node
	node2 := body.Nodes[1]

	fmt.Println("one-", node1, "two-",node2)

	n1, err1 := db.FindNode(node1.ID)
	n2, err2 := db.FindNode(node2.ID)

	if err1 != nil { //Node 1 is not in database
		n1, err1  = db.AddNode(&node1)

		if err1 != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while adding node1 to database",})
			return
		}
	}	


	if err2 != nil { //Node 2 is not in database
		n2, err2  = db.AddNode(&node2)

		if err2 != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while adding node2 to database",})
			return
		}
	}

	//	Now both nodes are guaranteed to be in the database

	edge, err := db.AddEdge(n1, n2, body.Dist) //the third item in body should be the cost of the edge traversal, or the real distance between n1 and n2

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while adding the edge to database",})
		return		
	}

	c.JSON(http.StatusOK, gin.H{"edge added": edge,})
	return
}

func allNodes(c *gin.Context){
	nodes, err := db.AllNodes()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
	return
}

func allEdges(c *gin.Context){
	edges, err := db.AllEdges()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"edges": edges,
	})
	return
}

func SetupRouter() *gin.Engine{
	r := gin.Default()


	r.POST("/addNode", addNode)
	r.POST("/addEdge", addEdge)
	r.GET("/allEdges", allEdges)
	r.GET("/allNodes", allNodes)

	return r
}