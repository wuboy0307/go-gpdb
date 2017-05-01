package main

import (
	"./src/download"
	"./src/install"
	"./src/argParser"
	"./src/core"
	log "./pkg/core/logger"
)


import (
	"./pkg/core/arguments"
	"./pkg/core/methods"
)

func main() {

	// Initialize the logger
	log.InitLogger()

	// Get all the configs
	core.Config()

	// Extract all the OS command line arguments
	if !arguments.InstallAfterDownload {
		argParser.ArgParser()
	}

	// Run Program based on what option is specified
	switch arguments.ArgOption {
	case "download":                                                // Run Download
		err := download.Download()
		methods.Fatal_handler(err)
		if arguments.InstallAfterDownload {
			arguments.ArgOption = "install"
			main()
		}
	case "install":                                                 // Run Install
		err := install.Install()
		methods.Fatal_handler(err)
	case "remove":                                                  // Run Remove
		log.Println("Run Remove")
	case "env":                                                     // Run env
		log.Println("Run env")
	default:                                                        // Error if command is invalid
		log.Fatal("Command not recognized")
	}

}