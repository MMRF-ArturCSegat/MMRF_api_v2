package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/sessions"
	"github.com/UFSM-Routelib/routelib_api/util"
	"github.com/gin-gonic/gin"
)


func parse_csv_to_obj(c * gin.Context){
    limiter_string := c.PostForm("limiter")
    fmt.Println("limiter: " + limiter_string)
    var coord_limiter util.Square
    err := json.Unmarshal([]byte(limiter_string), &coord_limiter)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"bad coord limiter": err.Error()})
        return
    }
    
    olt_string := c.PostForm("OLT")
    fmt.Printf("olt: %v\n", olt_string)
    var olt util.Coord
    err = json.Unmarshal([]byte(olt_string), &olt)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"bad olt cord: ": err.Error()})
        return
    }


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
    
    csvg, err := graph_model.New_csvg(file, olt, coord_limiter)
    file.Close()
    if err != nil{
        println("bad request")
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
        c.JSON(http.StatusOK, gin.H{"drawablePaths": csvg.Csvg_to_slice_of_coord_paths()})
        sessions.PrintSessions()
        return
    }
    /// if here, the client making the request already has a session logged in
    sessions.PrintSessions()
    c.JSON(http.StatusUnauthorized, gin.H{"error": "You already have a session created, delete it do create a new one"})
}

func has_session(c * gin.Context){
    cookie_string, err := c.Cookie("session_id")
    fmt.Println("searhing for " + cookie_string + "in cookie, jar")
    cookie, cookie_err := sessions.GetServerCookie(cookie_string)
    if err != nil || cookie_err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"session": "no session"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"session": cookie})
}


func delete_session(c * gin.Context){
    cookie_string, err := c.Cookie("session_id")
    cookie_err := sessions.RemoveSession(cookie_string)
    if err != nil || cookie_err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"session": "no session to delete"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"session": "deleted session" + cookie_string})
}


func generate_drawable_paths(c * gin.Context){
    cookie_string, err := c.Cookie("session_id")
    _, cookie_err := sessions.GetServerCookie(cookie_string)
    csvg, csvg_err := sessions.GetCSVG(cookie_string)
    if err != nil || cookie_err != nil || csvg_err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"session": "invalid cookie or expired session"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"drawablePaths": csvg.Csvg_to_slice_of_coord_paths()})
}
