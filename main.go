package main

import (
	"github.com/bujosa/xihe/scripts"
	"github.com/bujosa/xihe/utils"
)

func main() {
	utils.SetLogFile("log.txt")
	// Data Transformation Pipeline
	// transformation.RunCarTransformation()

	scripts.UploadDealers()

	// scripts.TrimMatchingStrategy()
}
