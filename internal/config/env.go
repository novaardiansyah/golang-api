package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	WebURL    string
	JWTSecret string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	WebURL = os.Getenv("WEB_URL")
	JWTSecret = os.Getenv("JWT_SECRET")
}
