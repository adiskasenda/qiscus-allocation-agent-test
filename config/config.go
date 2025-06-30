// package config

// import "os"

// var (
// 	BaseUrl   = os.Getenv("BASE_URL")
// 	SecretKey = os.Getenv("SECRET_KEY")
// 	AppCode   = os.Getenv("APP_ID")
// )

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BaseUrl    string
	SecretKey  string
	AppCode    string
	dbUser     string
	dbPassword string
	dbName     string
	dbHost     string
	dbDriver   string
	dbPort     string
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file not found, continuing with system env")
	}

	BaseUrl = os.Getenv("BASE_URL")
	SecretKey = os.Getenv("SECRET_KEY")
	AppCode = os.Getenv("APP_ID")

	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
	dbDriver = os.Getenv("DB_DRIVER")
	dbPort = os.Getenv("DB_PORT")
}
