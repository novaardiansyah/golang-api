package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	WebURL    string
	JWTSecret string

	MailHost        string
	MailPort        int
	MailUsername    string
	MailPassword    string
	MailEncryption  string
	MailFromAddress string
	MailFromName    string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	WebURL = os.Getenv("WEB_URL")
	JWTSecret = os.Getenv("JWT_SECRET")

	MailHost = os.Getenv("MAIL_HOST")
	MailPort, _ = strconv.Atoi(os.Getenv("MAIL_PORT"))
	MailUsername = os.Getenv("MAIL_USERNAME")
	MailPassword = os.Getenv("MAIL_PASSWORD")
	MailEncryption = os.Getenv("MAIL_ENCRYPTION")
	MailFromAddress = os.Getenv("MAIL_FROM_ADDRESS")
	MailFromName = os.Getenv("MAIL_FROM_NAME")
}
