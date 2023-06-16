package routes

import (
	"fmt"
	"net/http"
	"os"
	gm "github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/sessions"
	"github.com/gin-gonic/gin"
	sgo"github.com/UFSM-Routelib/routelib_api/sub_graph_optimization"
)

func generate_txt(c * gin.Context){
    type Body struct {
        Paths          []gm.GraphPath       `json:"paths"`
        Cables         []uint               `json:"cables"`
        Spliceboxes    []uint               `json:"boxes"`
        Uspliters      []uint               `json:"uspliters"`
        Bspliters      []uint               `json:"bspliters"`
    }

    var body Body
	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

    cookie_string, cookie_err := c.Cookie("session_id")
    if cookie_err != nil{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no valid session"})
        return 
    }
    if _, err := sessions.GetServerCookie(cookie_string); err != nil{
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return 
    }

    sub_graph := gm.Slice_of_paths_to_csvg(body.Paths)
    file, file_err :=  sgo.GenerateSubGraphOptimizationFile(sub_graph, body.Cables, body.Spliceboxes, body.Uspliters, body.Bspliters)
    if file_err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": file_err.Error()})
        return
    }
    defer file.Close()

    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "sub_graph.txt"))
    c.Header("Content-Type", "text/plain")
    c.Header("Content-Description", "File Transfer")
    c.File("/home/arturcs/Documents/routelib_api/sub_graph.txt")
    os.Remove("/home/arturcs/Documents/routelib_api/sub_graph.txt")
}





