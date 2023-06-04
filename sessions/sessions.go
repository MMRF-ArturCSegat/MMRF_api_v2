package sessions

import (
	"errors"
	"fmt"
	"time"

	gp "github.com/UFSM-Routelib/routelib_api/graph_model"
)

// global variable for the server to store its graphs, maped to its respective client creator
var Sessions = make(map[ServerCookie]*gp.CSV_Graph)
var cookie_jar = make(map[string]ServerCookie)

func AddSession(cookie_in ServerCookie, csvg *gp.CSV_Graph) error {
    _, exists := cookie_jar[cookie_in.ID]
    if exists {
        return errors.New("failed to add, this client already has a session")
    }
    cookie_jar[cookie_in.ID] = cookie_in
    Sessions[cookie_in] = csvg
    return nil
}


func RemoveSession(cookie_id string) error {
    cookie, exists := cookie_jar[cookie_id]
    if exists{
        delete(Sessions, cookie)
        delete(cookie_jar, cookie.ID)
        return nil
    }
    return errors.New("This session does not exist, thus cannot be deleted")
} 


func GetServerCookie(id string) (ServerCookie, error){
    cookie, exists := cookie_jar[id]
    if exists{
        if cookie.isExpired(){
            delete(cookie_jar, cookie.ID)
            return ServerCookie{ID:"", expiry: time.Now()}, errors.New("this cookie is expired")
        }
        return cookie, nil
    }
    return ServerCookie{ID:"", expiry: time.Now()}, errors.New("cookie does not exist")
}


func GetCSVG(cookie_id string) (*gp.CSV_Graph, error){
    cookie, exists := cookie_jar[cookie_id]

    if exists{
        if cookie.isExpired(){
            delete(Sessions, cookie)
            delete(cookie_jar, cookie.ID)
            return nil, errors.New("this session has expired")
        }
        return Sessions[cookie], nil
    }
    return nil, errors.New("No session registered")
}


func CleanExpiredSessions(){
    for key := range Sessions{
        if key.isExpired(){
            delete(Sessions, key)
            delete(cookie_jar, key.ID)
        }
    }
}


func PrintSessions(){
    for cookie, csvg := range Sessions{
        fmt.Printf("%v: \n", cookie.ID)
        csvg.Print()
    }
}
