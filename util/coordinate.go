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

type Distance float32

// Constants needed for distance calculations
const (
	EarthRadius       = 6371 * 1000.0
	DoubleEarthRadius = 2 * EarthRadius
	PiOver180         = math.Pi / 180
)

// DistanceBetween calculates the distance between two coordinates
func DistanceBetween(a, b Coord) float32 {
	value := 0.5 - math.Cos((b.Lat-a.Lat)*PiOver180)/2 + math.Cos(a.Lat*PiOver180)*math.Cos(b.Lat*PiOver180)*(1-math.Cos((b.Lng-a.Lng)*PiOver180))/2
	return DoubleEarthRadius * float32(math.Asin(math.Sqrt(value)))
}

// DistanceTo calculates the distance from this coordinate to another coordinate
func (c Coord) DistanceToInMeters(other Coord) float32 {
    dist :=  DistanceBetween(c, other)
	return dist
} 


func (co *Coord) IsInSquare(sq Square) bool {
    validLat := co.Lat <= sq.Top.Lat && co.Lat >= sq.Bot.Lat
    validLng := co.Lng >= sq.Top.Lng && co.Lng <= sq.Bot.Lng

    return validLat && validLng
}

