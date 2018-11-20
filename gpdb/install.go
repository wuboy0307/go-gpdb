package main

import (
	"fmt"
	"os"
	"strconv"
)

func install() {

	Infof("Running the installation for the product: %s", cmdOptions.Product)

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
	unzip(fmt.Sprintf("*%s*.zip", cmdOptions.Version))

}


// Install GPCC
func installGPCC() {


}

