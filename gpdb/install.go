package main

import (
	"os"
)

type Installation struct {
	HostFileLocation string
	WorkingHostFileLocation string
	SegmentHostLocation		string
	SegInstallHostLocation	string
	BinaryInstallationLocation string
	GpInitSystemConfigLocation string
	PortFileName string
	EnvFile string
	Timestamp string
	SingleORMulti string
	StandbyHostAvailable bool
	GPInitSystem GPInitSystemConfig
	GPCC GPCCConfig
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

type GPCCConfig struct {
	InstanceName string
	InstancePort string
	GpPerfmonHome string
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
		i.installGPCC()
	}
}

// Install GPDB
func (i *Installation) installGPDB(singleOrMutli string) {

	Infof("Starting the program to install GPDB version: %s", cmdOptions.Version)

	// Pre check
	i.preGPDBChecks()

	// Install the product
	i.installGPDBProduct()

	// run post steps
	i.postGPDBInstall()

	Infof("Installation of GPDB with version %s is complete", cmdOptions.Version)
	Infof("exiting ....")
}


// Install GPCC
func (i *Installation) installGPCC() {

	Infof("Starting the program to install GPCC version \"%s\" on the GPDB Version \"%s\"", cmdOptions.CCVersion, cmdOptions.Version)

	// Run prechecks
	i.preGPCCChecks()

	// Install GPCC product
	i.installGPCCProduct()

	// run post steps
	i.postGPCCInstall()

	Infof("Installation of GPCC with version %s on GPDB with version %s is complete", cmdOptions.CCVersion, cmdOptions.Version)
	Infof("exiting ....")
}

