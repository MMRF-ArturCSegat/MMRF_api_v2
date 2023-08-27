package routes

import (
	"net/http"
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
    r.GET("/has-session/", has_session)
    r.GET("/delete-session/", delete_session)

    // graph interaction
    r.POST("/upload_csv/", parse_csv_to_obj)
	r.GET("/all-nodes/", allNodes)
    r.POST("/limited-branching/", SpreadRadius)
    r.POST("/closest-node/", ClosestNode)
    r.POST("/txt-sub-graph/", generate_txt)
    r.GET("/drawable-paths/", generate_drawable_paths)

    // fiber optic component storage
    r.GET("/get-all-cables/", all_cables)
    r.GET("/get-cable/:id/", get_cable)
    r.POST("/add-cable/", add_cable)

    r.GET("/get-all-spliceboxes/", all_boxes)
    r.GET("/get-splicebox/:id/", get_box)
    r.POST("/add-splicebox/", add_box)

    r.GET("/get-all-uspliters/", all_u_spliters)
    r.GET("/get-uspliter/:id/", get_u_spliter)
    r.POST("/add-uspliter/", add_u_spliter)

    r.GET("/get-all-bspliters/", all_b_spliters)
    r.GET("/get-bspliter/:id/", get_b_spliter)
    r.POST("/add-bspliter/", add_b_spliter)
	return r
}
