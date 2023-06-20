package routes

import (
	"fmt"
	"net/http"

	gm"github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/sessions"
	"github.com/UFSM-Routelib/routelib_api/util"
	"github.com/gin-gonic/gin"
)

func allNodes(c *gin.Context){
    cookie_string, cookie_err := c.Cookie("session_id")
    csvg, err := sessions.GetCSVG(cookie_string)
    if err != nil || cookie_err != nil{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no valid session"})
        return 
    }

    nodes := csvg.AllNodes()

    c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}


func SpreadRadius(c * gin.Context){

    type Body struct{
        Node        gm.GraphNode    `json:"node"`
        Cost        float32         `json:"cost"`
        Limit       float32         `json:"limit"`
        Square      util.Square     `json:"square"`
    }
   
    var body Body
   
	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

    cookie_string, cookie_err := c.Cookie("session_id")
    csvg, csvg_err := sessions.GetCSVG(cookie_string)
    if cookie_err != nil || csvg_err !=  nil{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no valid session"})
    }


    start, err := csvg.FindNode(body.Node.ID)
    fmt.Println("bn", body)

    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return 
    }
                                //  this subtracion happens so a inital cost can be taken into account
    paths := csvg.LimitedBranchigFrom(start, (body.Limit - body.Cost),  gm.GraphPath{Nodes: make([]*gm.GraphNode, 0), Cost: 0}, make([]gm.GraphPath, 0), body.Square)

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

    cookie_string, cookie_err := c.Cookie("session_id")
    csvg, csvg_err := sessions.GetCSVG(cookie_string)
    if cookie_err != nil || csvg_err !=  nil{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no valid session"})
    }

   
    node, dist := csvg.ClosestNode(coord)
    println(node, dist)
   
    type result struct{
        Node    *gm.GraphNode   `json:"node"`
        Dist    float32         `json:"dist"`
    }

    c.JSON(http.StatusOK, gin.H{"closest-pair": result{Node: node, Dist: dist},})
}
