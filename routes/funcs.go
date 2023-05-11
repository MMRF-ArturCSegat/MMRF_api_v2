package routes

import (
	"fmt"
	"github.com/MMRF-ArturCSegat/MMRF_api_v2/db"
    "github.com/MMRF-ArturCSegat/MMRF_api_v2/util"
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
        Cost        float32         `json:"cost"`
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
                                //  this subtracion happens so a inital cost can be taken into account
    paths := db.SpreadRadius(start, (body.Limit - body.Cost),  db.GraphPath{Nodes: make([]*db.Node, 0), Cost: 0}, make([]db.GraphPath, 0), body.Square)
	
	for _, path := range paths{
        fmt.Println(path.IdSlice(), "cost: ", path.Cost)
	}

    c.JSON(http.StatusOK, gin.H{"paths":paths,})
}


func ClosestNode(c * gin.Context){
    
    var coord util.Coord

    if err := c.ShouldBindJSON(&coord); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
    }
    
    node, dist := db.ClosestNode(coord)
    println(node, dist)
    
    type result struct{
        Node *db.Node `json:"node"`
        Dist float32  `json:"dist"`
    }

    c.JSON(http.StatusOK, gin.H{"closest-pair": result{Node: node, Dist: dist},})
}
