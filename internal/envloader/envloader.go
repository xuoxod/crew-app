package envloader

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVars() error {
	// Load environment variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	return nil
}
