package main

import (
	"fmt"
	"os"
	"strconv"
)

func install() {

	Infof("Running the installation for the product: %s", cmdOptions.Product)

	// Get the hostname
	getHostname()

	// Checking if this is a single install VM or Mutli node VM
	var singleORmulti string
	noSegments, _ := strconv.Atoi(os.Getenv("GPDB_SEGMENT"))
	if noSegments > 0 {
		singleORmulti = "multi"
	} else {
		singleORmulti = "single"
	}
	Debugf("Is this single or multi node installation: %s", singleORmulti)

	// Run the installation
	if cmdOptions.Product == "gpdb" { // Install GPDB
		installGPDB(singleORmulti)
	} else { // its a GPCC installation
		installGPCC()
	}

}

// Install GPDB
func installGPDB(singleOrMutli string) {

	Infof("Starting the program to install GPDB version: %s", cmdOptions.Version)

	// Validate the master & segment exists and is readable
	dirValidator()

	// TODO: Check if there is already a version of GPDB installed

	// Check if the binaries exits and unzip the binaries.
	binFile := unzip(fmt.Sprintf("*%s*", cmdOptions.Version))

	Infof("Using the bin file to install the GPDB Product: %s", binFile)

	Infof("Installation of GPDB with version %s is complete", cmdOptions.Version)

}


// Install GPCC
func installGPCC() {

	Infof("Installation of GPCC with version %s on GPDB with version %s is complete", cmdOptions.CCVersion, cmdOptions.Version)
}

