package db

import (
	pq "github.com/lib/pq"
)


type Node struct {
	ID		int64		`json:"id" gorm:"primaryKey"`
	Lat		float64		`json:"lat"`
	Lng		float64		`json:"lng"`
	Neighbours	pq.Int64Array	`json:":neighbours" gorm:"type:integer[]"`
}










