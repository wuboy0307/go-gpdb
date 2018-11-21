package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Installation struct {
	HostFileLocation string
	BinaryInstallationLocation string
	SingleORMulti string
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
	i.getHostname()

	// Run the installation
	if cmdOptions.Product == "gpdb" { // Install GPDB
		i.installGPDB(singleORmulti)
	} else { // its a GPCC installation
		installGPCC()
	}

}

// Get and generate host file if doesn't exists the hostname
func (i *Installation) getHostname() {

	i.HostFileLocation = fmt.Sprintf("%s/hostfile", os.Getenv("HOME"))

	exists, _ := doesFileOrDirExists(i.HostFileLocation)
	if !exists {
		Infof("Host file doesn't exists, creating one: %s", i.HostFileLocation)
		// Read the contents from the /etc/hosts and generate a hostfile.
		hosts := contentExtractor(readFile("/etc/hosts"), "{if (NR!=1) {print $2}}", []string{})

		// Replace the last blank lines
		var s = hosts.String()
		regex, err := regexp.Compile("\n$")
		if err != nil {
			return
		}
		s = regex.ReplaceAllString(s, "")

		// write to the file
		writeFile(i.HostFileLocation, []string{s})
	} else {
		Infof("Found host file at location: %s", i.HostFileLocation)
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
	err := CheckHostnameIsValid()
	if err != nil {
		Fatalf("")
	}

	Infof("Installation of GPDB with version %s is complete", cmdOptions.Version)

}


// Install GPCC
func installGPCC() {

	Infof("Installation of GPCC with version %s on GPDB with version %s is complete", cmdOptions.CCVersion, cmdOptions.Version)
}

