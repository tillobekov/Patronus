package util

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GoDotEnvVariable(key string) string {
	err := godotenv.Load("app.env")

	if err != nil {
		log.Fatalf("Error loading app.env file")
	}

	return os.Getenv(key)
}
