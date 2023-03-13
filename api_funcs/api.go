package funcs

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"net/http"
	"gat/test"
	"fmt"
)

func home(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message": "Lol funny api",
	})
	return
}

func addNode(c *gin.Context){
	var node db2.Node

	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in binding": err.Error(),})
		return
	}

	fmt.Println(node)

	res, err := db2.AddNode(&node)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error in database:": err.Error(),})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"node added": res})
	return
}

func addShit(c  *gin.Context){
	nid ,_ := strconv.Atoi(c.Param("id"))
	node, err := db2.AddShit(nid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to add shit",})
		return
	}

	c.JSON(http.StatusTeapot, gin.H{"shit": node,})
	return
}

func allNodes(c *gin.Context){
	nodes, err := db2.AllNodes()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
	return
}

func addEdge(c *gin.Context){

	type Body struct {
		Nodes 	[]db2.Node 	`jason:"nodes"`
		Dist  	int		`jason:"dist"`
	}

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	if len(body.Nodes) != 2{
		c.JSON(http.StatusBadRequest, gin.H{"error": "expecting an array of 2 user objects"})
        return
	}

	node1 := body.Nodes[0] // start and end node
	node2 := body.Nodes[1]

	fmt.Println("one-", node1, "two-",node2)
	nodes, err := db2.AddEdge(&node1, &node2)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes" : nodes,})
	return
}

func connect(c * gin.Context){

	type Body struct {
		Nodes 	[]db2.Node 	`jason:"nodes"`
		Dist  	int		`jason:"dist"`	
	}

	var nodes Body

	if err := c.ShouldBindJSON(&nodes); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}
	err := db2.ConnectNodes(&nodes.Nodes[0], &nodes.Nodes[1])
	
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes" : nodes,})
	return
}

func SetupRouter() *gin.Engine{
	r := gin.Default()


	r.POST("/addNode", addNode)
	r.POST("/addedge", addEdge)
	// r.get("/alledges", alledges)
	r.GET("/allNodes", allNodes)
	r.GET("addShit/:id", addShit)
	r.PUT("/addEdge", connect)
	return r
}
