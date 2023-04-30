package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func home(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message": "Lol funny api",
	})
}

func SetupRouter() *gin.Engine{
	r := gin.Default()
	r.GET("/", home)
	r.POST("/add-node", addNode)
	r.POST("/add-edge", addEdge)
	r.GET("/all-nodes", allNodes)
	r.POST("/connect", connect)
	r.POST("/spread-radius/", SpreadRadius)
	return r
}
