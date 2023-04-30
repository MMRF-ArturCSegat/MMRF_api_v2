package routes

import (
	"fmt"
	"gat/db"
    "gat/util"
	"net/http"
	"github.com/gin-gonic/gin"
)

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
}


func SpreadRadius(c * gin.Context){

    type Body struct{
        node_id         int64
        limit           float32
        square          util.Square
    }
    
    var body Body
    
	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

    start, err := db.FindNode(body.node_id)

    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return 
    }

    paths := db.SpreadRadius(start, body.limit,  db.GraphPath{Nodes: make([]*db.Node, 0), Cost: 0}, make([]db.GraphPath, 0), body.square)
	
	for _, path := range paths{
        fmt.Println(path.IdSlice(), "cost: ", path.Cost)
	}

    c.JSON(http.StatusOK, gin.H{"paths":paths,})
}
