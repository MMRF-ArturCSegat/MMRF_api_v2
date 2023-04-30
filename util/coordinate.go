package util

import (
    "math"
)

type Coord struct {
    lat float64
    lng float64
}

type Square struct{
    top Coord
    bot Coord
}


func NewCoord(lat, lng float64) Coord {
    return Coord {lat: lat, lng: lng}
} 

func (co * Coord) DistanceToInMeters(dest Coord) float32 {
    lat1 := co.lat
    lon1 := co.lng
    lat2 := dest.lat
    lon2 := dest.lng
    R := 6371000.0 // Earth radius in meters
    phi1 := lat1 * math.Pi / 180.0
    phi2 := lat2 * math.Pi / 180.0
    deltaPhi := (lat2 - lat1) * math.Pi / 180.0
    deltaLambda := (lon2 - lon1) * math.Pi / 180.0
    a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    distance := R * c
    return float32(distance)
}

func (co * Coord) IsInSquare(sq Square)  bool{
    valid_lat := false
    valid_lng := false

    if co.lat >= sq.bot.lat && co.lat <= sq.top.lat{
        valid_lat = true
    }

    if co.lng >= sq.bot.lng && co.lng <= sq.top.lng{
        valid_lng = true
    }

    if valid_lng && valid_lat{
        return true
    }
    
    return false 
}
