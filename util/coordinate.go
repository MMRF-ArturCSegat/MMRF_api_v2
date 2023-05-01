package util

import (
    "math"
)

type Coord struct {
    Lat         float64     `json:"lat"`
    Lng         float64     `json:"lng"`
}

type Square struct{
    Top         Coord       `json:"top"`
    Bot         Coord       `json:"bot"`
}


func NewCoord(lat, lng float64) Coord {
    return Coord {Lat: lat, Lng: lng}
} 

func (co * Coord) DistanceToInMeters(dest Coord) float32 {
    lat1 := co.Lat
    lon1 := co.Lng
    lat2 := dest.Lat
    lon2 := dest.Lng
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

    if co.Lat >= sq.Bot.Lat && co.Lat <= sq.Top.Lat{
        valid_lat = true
    }

    if co.Lng >= sq.Bot.Lng && co.Lng <= sq.Top.Lng{
        valid_lng = true
    }

    if valid_lng && valid_lat{
        return true
    }
    
    return false 
}
