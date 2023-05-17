package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB



func loadEnvVariables(filePath string) (map[string]string, error) {
	err := godotenv.Load(filePath)
	if err != nil {
		return nil, err
	}


    getEnvVariable := func (key string) string {
        value := os.Getenv(key)
        return value
    }

	envVars := make(map[string]string)
	envVars["DB_HOST"] = getEnvVariable("DB_HOST")
	envVars["DB_USER_NAME"] = getEnvVariable("DB_USER_NAME")
	envVars["DB_PSSW"] = getEnvVariable("DB_PSSW")
	envVars["DB_NAME"] = getEnvVariable("DB_NAME")
	envVars["DB_PORT"] = getEnvVariable("DB_PORT")

	return envVars, nil
}



func ConnectDatabase2() {
    envVars, err := loadEnvVariables(".env")

    if err != nil{
        fmt.Println("Failed to load .env vairables")
        return
    }

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
                    envVars["DB_HOST"],envVars["DB_USER_NAME"],envVars["DB_PSSW"],envVars["DB_NAME"],envVars["DB_PORT"], )
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
