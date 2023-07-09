package main

import (
	"github.com/bujosa/xihe/transformation"
	"github.com/bujosa/xihe/utils"
)

func main() {
	utils.SetLogFile("log.txt")
	// Data Transformation Pipeline
	transformation.RunCarTransformation()

	// scripts.TrimMatchingStrategy()
}
