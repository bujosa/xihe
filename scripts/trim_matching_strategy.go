package scripts

import (
	"log"
	"os"

	"github.com/bujosa/xihe/database"
)

func TrimMatchingStrategy() {
	// Add register file
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.SetPrefix("[INFO] ")
	log.Println("Starting Trim Matching Strategy...")

	cars := database.GetCars()

	print(len(cars))
}