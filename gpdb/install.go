package main

import (
	"os"
	"strconv"
)

type Installation struct {
	HostFileLocation string
	WorkingHostFileLocation string
	BinaryInstallationLocation string
	SingleORMulti string
	portFileName string
	timestamp string
	GPInitSystem GPInitSystemConfig
}

type GPInitSystemConfig struct {
	MasterHostname string
}

func install() {

	Infof("Running the installation for the product: %s", cmdOptions.Product)

	// Initialize the struct
	i := new(Installation)

	// Checking if this is a single install VM or Mutli node VM
	var singleORmulti string
	noSegments, _ := strconv.Atoi(os.Getenv("GPDB_SEGMENT"))
	if noSegments > 0 {
		i.SingleORMulti = "multi"
	} else {
		i.SingleORMulti = "single"
	}
	Debugf("Is this single or multi node installation: %s", singleORmulti)

	// Get or Generate the hostname file
	i.generateHostFile()

	// Run the installation
	if cmdOptions.Product == "gpdb" { // Install GPDB
		i.portFileName = "gpdb_ports.save"
		i.installGPDB(singleORmulti)
	} else { // its a GPCC installation
		i.portFileName = "gpcc_ports.save"
		installGPCC()
	}
}

// Install GPDB
func (i *Installation) installGPDB(singleOrMutli string) {

	Infof("Starting the program to install GPDB version: %s", cmdOptions.Version)

	// Validate the master & segment exists and is readable
	dirValidator()

	// TODO: Check if there is already a version of GPDB installed

	// Start the installation procedure
	i.installProduct()

	// Check ssh to host is working and enable password less login
	i.setUpHost()

	Infof("Installation of GPDB with version %s is complete", cmdOptions.Version)
}


// Install GPCC
func installGPCC() {

	Infof("Installation of GPCC with version %s on GPDB with version %s is complete", cmdOptions.CCVersion, cmdOptions.Version)
}

