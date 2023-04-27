package routes

import(
    "fmt"
    "strconv"
    "gat/db"
    "net/http"
	"github.com/gin-gonic/gin"
)

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
