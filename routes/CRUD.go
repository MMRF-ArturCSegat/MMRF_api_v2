package routes

import (
	"fmt"
	"net/http"

	"github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/sessions"
	"github.com/gin-gonic/gin"
)


func parse_csv_to_obj(c * gin.Context){
    file_ptr, file_err := c.FormFile("rede")
    if file_err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": file_err.Error()})
        println(file_err.Error())
        return
    }

    file, err := file_ptr.Open()
    if err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
        println(err.Error())
        return
    }
    
    csvg, err := graph_model.New_csvg(file)
    file.Close()
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        println(err.Error())
        return
    }
    //handle sessioning
    cookie_string, err := c.Cookie("session_id")
    fmt.Println("searhing for " + cookie_string + "in cookie, jar")
    _, cookie_err := sessions.GetServerCookie(cookie_string)
    if err != nil || cookie_err != nil{
        cookie := sessions.NewServerCookie()
        err = sessions.AddSession(cookie, csvg)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        // c.SetCookie has its MaxAge, Path, and Domain values zeroed because they are said to be optional in
        // https://pkg.go.dev/net/http#Cookie  AND https://github.com/gin-gonic/gin/blob/v1.9.0/context.go#L887
        // you may see i am not the most expireienced with cookie auth in gin
        c.SetCookie("session_id", cookie.ID, 0, "", "", false, false)
        c.JSON(http.StatusOK, gin.H{"message_string": "session created"})
        sessions.PrintSessions()
        return
    }
    /// if here, the client making the request already has a session logged in
    sessions.PrintSessions()
    c.JSON(http.StatusUnauthorized, gin.H{"error": "You already have a session created, delete it do create a new one"})
}


