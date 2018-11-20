package main

import "fmt"

func install() {

	Debugf("Running the installation for the product: %s", cmdOptions.Product)

	if cmdOptions.Product == "gpdb" { // Install GPDB
		installGPDB()
	} else { // its a GPCC installation
		installGPCC()
	}

}

// Install GPDB
func installGPDB() {

	Infof("Starting the program to install GPDB version: %s", cmdOptions.Version)

	// Validate the master & segment exists and is readable
	dirValidator()

	// TODO: Check if there is already a version of GPDB installed

	// Check if the binaries exits and unzip the binaries.
	unzip(fmt.Sprintf("*%s*.zip", cmdOptions.Version))

}


// Install GPCC
func installGPCC() {


}

