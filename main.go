package main

import (
	"flag"

	"github.com/bujosa/xihe/scripts"
	"github.com/bujosa/xihe/transformation"
	"github.com/bujosa/xihe/utils"
)

func main() {
	utils.SetLogFile("log.txt")
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
