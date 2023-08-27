package routes

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	ig"github.com/UFSM-Routelib/routelib_api/instance_generation"
	"github.com/UFSM-Routelib/routelib_api/sessions"
)

func generate_txt(c * gin.Context){

    var body ig.Instance
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
    
    file, file_err := body.GenerateSubGraphOptimizationFile(csvg)
    if file_err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": file_err.Error()})
        return
    }
    defer file.Close()

    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "sub_graph.txt"))
    c.Header("Content-Type", "text/plain")
    c.Header("Content-Description", "File Transfer")
    c.File("/home/arturcs/Documents/Routelib/routelib_api/sub_graph.txt")
    os.Remove("/home/arturcs/Documents/Routelib/routelib_api/sub_graph.txt")
}
