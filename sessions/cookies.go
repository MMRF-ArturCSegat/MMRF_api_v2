package sessions

import (
    "time"
    "github.com/google/uuid"
)


type ServerCookie struct {
    ID string
    expiry time.Time
}

func (sc * ServerCookie) isExpired() bool{
    return sc.expiry.Before(time.Now())
}

func NewServerCookie() ServerCookie {
    return ServerCookie{ID: uuid.NewString(), expiry: time.Now().Add(15 * time.Minute)}
}
