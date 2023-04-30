package util

import (
    "math"
)


func In[T comparable](e T, list []T) bool{
    for _, item := range list{
        if item == e{
            return true
        }
    }
    return false
}

func SliceComp[T comparable](s1, s2[]T) bool{
    for i, e := range s1{
        if s2[i] != e{
            return false
        }
    }
    return true
}

func SliceInMatrix[T comparable](m [][]T, s[]T) bool{
    for _, e := range m{
        if SliceComp(e, s) == true{
            return true
        }
    }
    return false
}

func CoordDistanceMeters(p1, p2 [2]float64) float64 {
    lat1 := p1[0]
    lon1 := p1[1]
    lat2 := p2[0]
    lon2 := p2[1]
    R := 6371000.0 // Earth radius in meters
    phi1 := lat1 * math.Pi / 180.0
    phi2 := lat2 * math.Pi / 180.0
    deltaPhi := (lat2 - lat1) * math.Pi / 180.0
    deltaLambda := (lon2 - lon1) * math.Pi / 180.0
    a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    distance := R * c
    return distance
}
