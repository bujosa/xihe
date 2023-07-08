package utils

import (
	"log"
	"os"
)

func SetLogFile(name string) {
	// Add register file
	logFile, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
}