package env

import (
	"fmt"
	"os"
)

type ErrEnvNotSet struct {
	key string
}

func (err *ErrEnvNotSet) Error() string {
	return fmt.Sprintf("Environment variable with key %s was not set", err.key)
}

// GetString reads a string value from the environment variables.
func GetString(key string) (string, error) {

	value, ok := os.LookupEnv(key)

	if !ok {
		return value, &ErrEnvNotSet{key: key}
	}

	return value, nil
}
