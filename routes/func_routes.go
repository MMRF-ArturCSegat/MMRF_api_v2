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
        Node        db.Node         `json:"node"`
        Limit       float32         `json:"limit"`
        Square      util.Square     `json:"square"`
    }
    
    var body Body
    
	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

    start, err := db.FindNode(body.Node.ID)
    fmt.Println("bn", body)

    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return 
    }

    paths := db.SpreadRadius(start, body.Limit,  db.GraphPath{Nodes: make([]*db.Node, 0), Cost: 0}, make([]db.GraphPath, 0), body.Square)
	
	for _, path := range paths{
        fmt.Println(path.IdSlice(), "cost: ", path.Cost)
	}

    c.JSON(http.StatusOK, gin.H{"paths":paths,})
}
