package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase2() {
    // Replace the connection details below with your own PostgreSQL database configuration
    dsn := "host=localhost user=arturcs password=123123123 dbname=gatdb port=5432 sslmode=disable TimeZone=UTC"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }

    err = database.AutoMigrate(&Node{})
    if err != nil {
        fmt.Println("Failed to auto migrate database schema")
        return
    }

    db = database
}
