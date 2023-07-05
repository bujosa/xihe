package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ErrEnvNotSet struct {
	key string
}

func (err *ErrEnvNotSet) Error() string {
	return fmt.Sprintf("Environment variable with key %s was not set", err.key)
}

// Get string from environment variable
func GetString(key string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	value := os.Getenv(key)

	if value == "" {
		return value, &ErrEnvNotSet{key: key}
	}

	return value, nil
}
