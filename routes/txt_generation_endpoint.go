package routes

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	gm "github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/sessions"
	sgo "github.com/UFSM-Routelib/routelib_api/sub_graph_optimization"
	"github.com/UFSM-Routelib/routelib_api/util"
	"github.com/gin-gonic/gin"
)

func generate_txt(c * gin.Context){
    cookie_string, cookie_err := c.Cookie("session_id")
    csvg, err := sessions.GetCSVG(cookie_string)
    if err != nil || cookie_err != nil{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no valid session"})
        return 
    }
    node_id := c.Param("node")
    parsed_node_id, err1 := strconv.ParseInt(node_id, 10, 64)
    node, err2 := csvg.FindNode(parsed_node_id)
    if err1 != nil || err2 != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    s := csvg.SpreadRadius(node, 999999999.0,  gm.GraphPath{Nodes: make([]*gm.GraphNode, 0), Cost: 0}, make([]gm.GraphPath, 0), util.DefaultMaxSquare())
    sub_graph := sgo.SubGraph(s)
    file, file_err := sub_graph.Build_txt_file()
    defer file.Close()
    if file_err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": file_err.Error()})
        return
    }

    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "sub_graph.txt"))
    c.Header("Content-Type", "text/plain")
    c.Header("Content-Description", "File Transfer")
    c.File("/home/arturcs/Documents/routelib_api/sub_graph.txt")
    os.Remove("/home/arturcs/Documents/routelib_api/sub_graph.txt")
}
