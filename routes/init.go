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
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
        c.Next()
    })
    r.GET("/", home)
    r.POST("/upload_csv/", parse_csv_to_obj)
	r.GET("/all-nodes/", allNodes)
    r.POST("/spread-radius/", SpreadRadius)
    r.POST("/closest-node/", ClosestNode)
    r.GET("/txt-sub-graph/:node", generate_txt)
	return r
}
