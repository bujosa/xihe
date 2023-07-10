package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/bujosa/xihe/scripts"
	"github.com/bujosa/xihe/transformation"
)

func main() {
	logName := "./logs/log_" + time.Now().Format("2006_01_02_15_04") + ".txt"
	logFile, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	var transformationCommand string
	var uploadCommand string
	var helpFlag bool

	flag.BoolVar(&helpFlag, "h", false, "Help")
	flag.StringVar(&transformationCommand, "t", "", "Transform data")
	flag.StringVar(&transformationCommand, "transform", "", "Transform data")
	flag.StringVar(&uploadCommand, "u", "", "Upload data")

	flag.Parse()

	if helpFlag {
		flag.PrintDefaults()
		return
	}

	if transformationCommand == "dealers" {
		transformation.RunDealerTransformation()
	} else if transformationCommand == "cars" {
		transformation.RunCarTransformation()
	}

	if uploadCommand == "dealers" {
		scripts.UploadDealers()
	} else if uploadCommand == "cars" {
		scripts.TrimMatchingStrategy()
	}
}
