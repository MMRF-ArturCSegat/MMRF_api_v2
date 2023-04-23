package routes

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"net/http"
	"gat/db"
	"fmt"
)

func home(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message": "Lol funny api",
	})
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
}

func allNodes(c *gin.Context){
	nodes, err := db.AllNodes()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
}

func addEdge(c *gin.Context){

	type Body struct {
		Nodes 	[]db.Node 	`jason:"nodes"`
		Dist  	int		`jason:"dist"`
	}

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	if len(body.Nodes) != 2{
		c.JSON(http.StatusBadRequest, gin.H{"error": "expecting an array of 2 node objects"})
        return
	}

	node1 := body.Nodes[0] // start and end node
	node2 := body.Nodes[1]

	fmt.Println("one-", node1, "two-",node2)
	nodes, err := db.AddEdge(&node1, &node2)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes" : nodes,})
}

func SpreadRadius(c * gin.Context){
	start_id, err := strconv.Atoi(c.Param("start"))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}
	
	start, err := db.FindNode(int64(start_id))
	
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	limit, err := strconv.Atoi(c.Param("limit"))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

    paths := db.SpreadRadius(start, limit, 0, make([]*db.Node, 0), make([][]*db.Node, 0))
	
	for _, e := range paths{
		fmt.Println(db.IdSliceFromNodeSlice(e))
	}

    c.JSON(http.StatusOK, gin.H{"paths":paths,})
}

func connect(c * gin.Context){
	n1, err1 := strconv.Atoi(c.Param("n1"))
	n2, err2 := strconv.Atoi(c.Param("n2"))

	if err1 != nil|| err2 != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error(),})
		return

	}

	err := db.ConnectNodes(int64(n1),int64(n2))
	
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node1" : n1, "node2": n2,})
}

func SetupRouter() *gin.Engine{
	r := gin.Default()
	r.GET("/", home)
	r.POST("/addNode", addNode)
	r.POST("/addEdge", addEdge)
	r.GET("/allNodes", allNodes)
	r.PUT("/addEdge/:n1/:n2", connect)
	r.GET("/spreadRadius/:start/:limit", SpreadRadius)
	return r
}
