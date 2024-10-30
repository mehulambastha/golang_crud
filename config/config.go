package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found. Loading with default values.")
	}
}
