package main

import (
	"./src/download"
	"./src/argParser"
	"./src/core"
	log "./pkg/core/logger"
)


func main() {

	// Initialize the logger
	log.InitLogger()

	// Get all the configs
	core.Config()

	// Extract all the OS command line arguments
	argParser.ArgParser()

	// Program to download the software from PivNet
	download.Download()
}