package routes

import (
	"net/http"
	// "github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

func home(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message": "Lol funny api",
	})
}





func SetupRouter() *gin.Engine{
	r := gin.Default()
    
	// Set up CORS middleware
    r.Use(func(c *gin.Context) {
	    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

        if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
        c.Next()
    })
    r.GET("/", home)
    r.POST("/add-node", addNode)
	r.POST("/add-edge", addEdge)
	r.GET("/all-nodes", allNodes)
	r.POST("/connect", connect)
	r.POST("/spread-radius/", SpreadRadius)
    r.POST("/closest-node/", ClosestNode)
	return r
}
