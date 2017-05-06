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
	err := core.Config()
	methods.Fatal_handler(err)

	// Extract all the OS command line arguments
	if !arguments.InstallAfterDownload {
		argParser.ArgParser()
	}

	// Run Program based on what option is specified
	switch arguments.ArgOption {
		case "download":
			// Run Download
			err := download.Download()
			methods.Fatal_handler(err)

			// If user has asked to run the installer after download then
			// set the argument and rerun the main to pick the right choice.
			if arguments.InstallAfterDownload {
				arguments.ArgOption = "install"
				arguments.RequestedInstallProduct = arguments.RequestedDownloadProduct
				main()
			}

		case "install":
			// If the product to install is GPDB then call the GPDB installation program
			if arguments.RequestedInstallProduct == "gpdb" { // Run Install GPDB
				err := install.InstallSingleNodeGPDB()
				methods.Fatal_handler(err)
			} else if arguments.RequestedInstallProduct == "gpcc" { // else GPCC installer
				err := install.InstalSingleNodeGPCC()
				methods.Fatal_handler(err)
			}
		case "remove":
			// Run remove
			log.Println("Run Remove")
		case "env":
			// Run env
			log.Println("Run env")
		default:
			// Error if command is invalid
			log.Fatal("Command not recognized")
	}

}