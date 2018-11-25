package main

import (
	"os"
	"time"
	"fmt"
)

type Installation struct {
	HostFileLocation string
	WorkingHostFileLocation string
	SegmentHostLocation		string
	BinaryInstallationLocation string
	SingleORMulti string
	PortFileName string
	Timestamp string
	StandbyHostAvailable bool
	GpInitSystemConfigLocation string
	GPInitSystem GPInitSystemConfig
	EnvFile string
}

type GPInitSystemConfig struct {
	MasterHostname string
	ArrayName	   string
	SegPrefix	   string
	DBName		   string
	MasterDir      string
	SegmentDir	   string
	MirrorDir	   string
	MasterPort	   string
	SegmentPort	   string
	MirrorPort	   string
	ReplicationPort string
	MirrorReplicationPort string
}

const (
	defaultMasterPort = 3000
	defaultGpperfmonPort = 28000
	defaultPrimaryPort = 30000
	defaultMirrorPort = 35000
	defaultReplicatePort = 40000
	defaultMirrorReplicatePort = 45000
)

func install() {

	Infof("Running the installation for the product: %s", cmdOptions.Product)

	// Initialize the struct
	i := new(Installation)

	// Checking if this is a single install VM or Mutli node VM
	var singleORmulti string
	noSegments := strToInt(os.Getenv("GPDB_SEGMENTS"))
	if noSegments > 0 {
		i.SingleORMulti = "multi"
	} else {
		i.SingleORMulti = "single"
	}
	Debugf("Is this single or multi node installation: %s", i.SingleORMulti)

	// Get or Generate the hostname file
	i.generateHostFile()

	// Run the installation
	if cmdOptions.Product == "gpdb" { // Install GPDB
		i.PortFileName = "gpdb_ports.save"
		i.installGPDB(singleORmulti)
	} else { // its a GPCC installation
		i.PortFileName = "gpcc_ports.save"
		installGPCC()
	}
}

// Install GPDB
func (i *Installation) installGPDB(singleOrMutli string) {

	Infof("Starting the program to install GPDB version: %s", cmdOptions.Version)

	// Validate the master & segment exists and is readable
	dirValidator()

	// Check if there is already a version of GPDB installed
	installedEnvFiles(fmt.Sprintf("*%s*", cmdOptions.Version), "confirm", true)

	// Start the installation procedure
	i.installProduct()

	// Check ssh to host is working and enable password less login
	i.setUpHost()

	// Build & Execute the gpinitsystem configuration
	i.Timestamp = time.Now().Format("20060102150405")
	i.buildGpInitSystem()

	// Create EnvFile
	i.createEnvFile()

	// Initialize standby master
	if i.StandbyHostAvailable && cmdOptions.Standby {
		i.activateStandby()
	} else if cmdOptions.Standby {
		Errorf("Cannot activate standby, please activate the standby manually")
	}

	// Create uninstall script
	i.createUninstallScript()

	// Store the last used port for future use
	i.savePort()

	// Installation complete, print on the screen the env file to source
	displayEnvFileToSource(i.EnvFile)
	defer deleteFile(i.WorkingHostFileLocation)
	defer deleteFile(i.SegmentHostLocation)

	Infof("Installation of GPDB with version %s is complete", cmdOptions.Version)
	Infof("exiting ....")
}


// Install GPCC
func installGPCC() {

	Infof("Installation of GPCC with version %s on GPDB with version %s is complete", cmdOptions.CCVersion, cmdOptions.Version)
	Infof("exiting ....")
}

