package fiber_optic_components

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectFBCDatabase(){
    database, err := gorm.Open(sqlite.Open("fiber_optic_components.db"), &gorm.Config{})

    if err != nil {
        panic("Failed to connect to database!")
    }

    err = database.AutoMigrate(&FiberCable{}, &FiberSpliceBox{}, &FiberBalancedSpliter{}, &FiberUnbalancedSpliter{})
    if err != nil {
        return
    }
    db = database
}

