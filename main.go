package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/bujosa/xihe/scripts"
	"github.com/bujosa/xihe/transformation"
	"github.com/bujosa/xihe/utils"
)

func main() {
	// Go routines
	runtime.GOMAXPROCS(1)

	// Create log file with timestamp
	logName := "./logs/log_" + time.Now().Format("2006_01_02_15_04") + ".txt"
	logFile, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// Load environment variables in context
	ctx := context.Background()
	utils.LoadEnvs(&ctx)

	// Define flags
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
		transformation.RunDealerTransformation(ctx)
	} else if transformationCommand == "cars" {
		transformation.RunCarTransformation(ctx)
	}

	if uploadCommand == "dealers" {
		scripts.UploadDealers(ctx)
	} else if uploadCommand == "cars" {
		scripts.TrimMatchingStrategy(ctx, false)
	} else if uploadCommand == "cars published" {
		scripts.TrimMatchingStrategy(ctx, true)
	}
}

func ReRunMatchingStrategy(ctx context.Context) {
	scripts.TrimMatchingStrategy(ctx, false)
}
