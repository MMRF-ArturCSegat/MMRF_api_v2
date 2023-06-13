package routes

import (
	"fmt"
	"net/http"
	"os"
	gm "github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/sessions"
	"github.com/gin-gonic/gin"
)

func generate_txt(c * gin.Context){
    var paths []gm.GraphPath
	if err := c.ShouldBindJSON(&paths); err != nil{
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

    sub_graph := gm.Slice_of_paths_to_csvg(paths)
    file, file_err := sub_graph.Build_txt_file()
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
