package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (i *Installation) buildGpInitSystem() {

	Infof("Building and executing the gpinitsystem...")

	// Set the values of the below parameters
	i.GPInitSystem.ArrayName = "gp_" + cmdOptions.Version + "_" + i.timestamp
	i.GPInitSystem.SegPrefix = "gp_" + cmdOptions.Version + "_" + i.timestamp
	i.GPInitSystem.DBName = "gpadmin"
	i.GPInitSystem.MasterDir = strings.TrimSuffix(Config.INSTALL.MASTERDATADIRECTORY, "/")
	i.GPInitSystem.SegmentDir = strings.TrimSuffix(Config.INSTALL.SEGMENTDATADIRECTORY, "/")
	i.GPInitSystem.MirrorDir = strings.TrimSuffix(Config.INSTALL.MIRRORDATADIRECTORY, "/")

	// Generate the port range
	i.generatePortRange()

	// Start building the gpinitsystem config file
	i.buildGpInitSystemConfig()

	// Stop all the database before execution
	stopAllDb()

	// Execute gpinitsystem
	i.executeGpInitSystem()

	// Is database healthy
	isDbHealthy("", i.GPInitSystem.MasterPort)
}

// Generate the port for master / segments / mirror & replication
func (i *Installation) generatePortRange() {

	Infof("Searching & Generating the port to be used for database installation")
	// Check if we have the last used port base file and its usable
	i.GPInitSystem.SegmentPort = i.validatePort( "PRIMARY_PORT", defaultPrimaryPort)  // segment
	i.GPInitSystem.MasterPort = i.validatePort("MASTER_PORT", defaultMasterPort) // master

	// If its a multi installation we will need the mirror / replication port as well & usable
	if i.SingleORMulti == "multi" {
		i.GPInitSystem.MirrorPort = i.validatePort("MIRROR_PORT", defaultMirrorPort) // mirror
		i.GPInitSystem.MirrorReplicationPort = i.validatePort("MIRROR_REPLICATION_PORT", defaultMirrorReplicatePort) // mirror replication port
		i.GPInitSystem.ReplicationPort = i.validatePort("REPLICATION_PORT", defaultReplicatePort) // replication
	}
}

// validate port
func (i *Installation) validatePort(searchString string, defaultPort int) string {
	Infof("Obtaining ports to be set for %s", searchString)
	p, _ := doWeHavePortBase(Config.INSTALL.FUTUREREFDIR, i.portFileName, searchString)
	if p == "" {
		Warnf("Didn't find %s in the file, setting it to default value: %d", searchString, defaultPort)
		p = strconv.Itoa(defaultPort)
	}
	p = strconv.Itoa(i.checkPortIsUsable(p))
	return p
}

// Building initsystem configuration
func (i *Installation) buildGpInitSystemConfig() {

	// Build gpinitsystem config file
	i.GpInitSystemConfigLocation = fmt.Sprintf("%s%s_%s_%s", Config.CORE.TEMPDIR, "gpinitsystemconfig", cmdOptions.Version, i.timestamp)
	Infof("Creating the gpinitsystem config file at: %s", i.GpInitSystemConfigLocation)
	deleteFile(i.GpInitSystemConfigLocation)
	createFile(i.GpInitSystemConfigLocation)

	// Write the below content to config file
	if i.SingleORMulti == "single" {
		writeFile(i.GpInitSystemConfigLocation, i.singleNodeGpInitSystem())
	} else {
		writeFile(i.GpInitSystemConfigLocation, i.multiNodeGpInitSystem())
	}
}

// The contents of single node gpinitsystem
func (i *Installation) singleNodeGpInitSystem() []string{
	Infof("Finalizing the gpinitsystem for the single mode database installation")
	return []string{
		"ARRAY_NAME=" + i.GPInitSystem.ArrayName,
		"SEG_PREFIX=" + i.GPInitSystem.SegPrefix,
		"MASTER_HOSTNAME=" + i.GPInitSystem.MasterHostname,
		"MASTER_DIRECTORY=" + i.GPInitSystem.MasterDir,
		"PORT_BASE=" + i.GPInitSystem.SegmentPort,
		"MASTER_PORT=" + i.GPInitSystem.MasterPort,
		"DATABASE_NAME=" + i.GPInitSystem.DBName,
		"declare -a DATA_DIRECTORY=("+ generateSegmentDirectoryList(i.GPInitSystem.SegmentDir) +")",
	}
}

// The contents of multi node gpinitsystem
func (i *Installation) multiNodeGpInitSystem() []string{
	Infof("Finalizing the gpinitsystem for the multi mode database installation")
	return []string{
		"ARRAY_NAME=" + i.GPInitSystem.ArrayName,
		"SEG_PREFIX=" + i.GPInitSystem.SegPrefix,
		"PORT_BASE=" + i.GPInitSystem.SegmentPort,
		"declare -a DATA_DIRECTORY=("+ generateSegmentDirectoryList(i.GPInitSystem.SegmentDir) +")",
		"MASTER_HOSTNAME=" + i.GPInitSystem.MasterHostname,
		"MASTER_DIRECTORY=" + i.GPInitSystem.MasterDir,
		"MASTER_PORT=" + i.GPInitSystem.MasterPort,
		"TRUSTED_SHELL=ssh",
		"CHECK_POINT_SEGMENTS=8",
		"ENCODING=UNICODE",
		"MIRROR_PORT_BASE=" + i.GPInitSystem.MirrorPort,
		"REPLICATION_PORT_BASE=" + i.GPInitSystem.ReplicationPort,
		"MIRROR_REPLICATION_PORT_BASE=" + i.GPInitSystem.MirrorReplicationPort,
		"declare -a MIRROR_DATA_DIRECTORY=("+ generateSegmentDirectoryList(i.GPInitSystem.MirrorDir) +")",
		"DATABASE_NAME=" + i.GPInitSystem.DBName,
	}
}

// Number of segment calculator
func generateSegmentDirectoryList(whichDir string) string {
	Debugf("Generating directory list of %s", whichDir)
	var dir string
	for i := 1; i <= Config.INSTALL.TOTALSEGMENT; i++ {
		dir = dir + " " + whichDir
	}
	return dir
}

// Check if the port is available or not
func (i *Installation) checkPortIsUsable(port string) int {
	Debugf("Checking for port \"%s\"is usable", port)
	pb := strToInt(port)
	pb, err := isPortUsed(pb, Config.INSTALL.TOTALSEGMENT, i.WorkingHostFileLocation)
	if err != nil {
		Fatalf("Error in checking the port usage, err: %v", err)
	}
	return pb
}

func (i *Installation) executeGpInitSystem() {
	Infof("Executing the gpinitsystem to install the database")
	executeOsCommand(fmt.Sprintf("%s/bin/gpinitsystem", os.Getenv("GPHOME")), "-c", i.GpInitSystemConfigLocation, "-h", i.SegmentHostLocation , "-a")
}