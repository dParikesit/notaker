package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Dotenv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Panicf("Error loading .env")
	}

	return os.Getenv(key)
}
