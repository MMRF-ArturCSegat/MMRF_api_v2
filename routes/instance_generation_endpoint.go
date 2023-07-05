package routes

import (
	"fmt"
	"net/http"
	"os"

	gm "github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/sessions"
	ig "github.com/UFSM-Routelib/routelib_api/instance_generation"
	"github.com/UFSM-Routelib/routelib_api/util"
	"github.com/gin-gonic/gin"
)

func generate_txt(c * gin.Context){
    type Body struct {
        Paths          [][]gm.GraphPath     `json:"paths"`
        Clients        []util.Coord         `json:"clients"`
        OLT            util.Coord           `json:"olt"`
        Cables         []uint32             `json:"cables"`
        Spliceboxes    []uint32             `json:"boxes"`
        Uspliters      []uint32             `json:"uspliters"`
        Bspliters      []uint32             `json:"bspliters"`
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
    csvg, csvg_err := sessions.GetCSVG(cookie_string)
    if cookie_err != nil || csvg_err !=  nil{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no valid session"})
    }
    
    var paths []*gm.CSV_Graph
    for _, path := range body.Paths{
        paths = append(paths, gm.Slice_of_paths_to_csvg(path))
    }
    file, file_err := ig.GenerateSubGraphOptimizationFile(csvg, paths, body.OLT, body.Clients, body.Cables, body.Spliceboxes, body.Uspliters, body.Bspliters)
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





