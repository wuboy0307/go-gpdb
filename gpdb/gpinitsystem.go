package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (i *Installation) buildGpInitSystem() {

	// Set the values of the below parameters
	i.GPInitSystem.ArrayName = "gp_" + cmdOptions.Version + "_" + i.timestamp
	i.GPInitSystem.SegPrefix = "gp_" + cmdOptions.Version + "_" + i.timestamp
	i.GPInitSystem.DBName = "gpadmin"
	i.GPInitSystem.MasterDir = strings.TrimSuffix(Config.INSTALL.MASTERDATADIRECTORY, "/")
	Config.INSTALL.SEGMENTDATADIRECTORY = strings.TrimSuffix(Config.INSTALL.SEGMENTDATADIRECTORY, "/")

	// Generate the port range
	i.generatePortRange()

	// Start building the gpinitsystem config file
	i.buildGpInitSystemConfig()

	// Execute gpinitsystem
	i.executeGpInitSystem()
}

func (i *Installation) generatePortRange() {

	// Check if we have the last used port base file and its usable
	Infof("Obtaining ports to be set for primary segments")
	i.GPInitSystem.SegmentPort, _ = doWeHavePortBase(Config.INSTALL.FUTUREREFDIR, i.portFileName, "PRIMARY_PORT")
	if i.GPInitSystem.SegmentPort == "" {
		Warnf("Didn't find PRIMARY_PORT in the file, setting it to default value: %d", defaultPrimaryPort)
		i.GPInitSystem.SegmentPort = strconv.Itoa(defaultPrimaryPort)
	}
	i.GPInitSystem.SegmentPort = strconv.Itoa(i.checkPortIsUsable(i.GPInitSystem.SegmentPort))

	// Check if we have the last used master port file and if its usable
	Infof("Obtaining ports to be set for master segment")
	i.GPInitSystem.MasterPort, _ = doWeHavePortBase(Config.INSTALL.FUTUREREFDIR, i.portFileName, "MASTER_PORT")
	if i.GPInitSystem.MasterPort == "" {
		Warnf("Didn't find MASTER_PORT in the file, setting it to default value: %d", defaultMasterPort)
		i.GPInitSystem.MasterPort = strconv.Itoa(defaultMasterPort)
	}
	i.GPInitSystem.MasterPort = strconv.Itoa(i.checkPortIsUsable(i.GPInitSystem.MasterPort))

	// If its a multi installation we will need the mirror port as well & usable
	if i.SingleORMulti == "multi" {
		// Check if we have the last used mirror port file
		Infof("Obtaining ports to be set for mirror segment")
		i.GPInitSystem.MirrorPort, _ = doWeHavePortBase(Config.INSTALL.FUTUREREFDIR, i.portFileName, "MIRROR_PORT")
		if i.GPInitSystem.MirrorPort == "" {
			Warnf("Didn't find MIRROR_PORT in the file, setting it to default value: %d", defaultMirrorPort)
			i.GPInitSystem.MirrorPort = strconv.Itoa(defaultMirrorPort)
		}
		i.GPInitSystem.MirrorPort = strconv.Itoa(i.checkPortIsUsable(i.GPInitSystem.MirrorPort))

		// Check if we have the last used replication port file
		Infof("Obtaining ports to be set for replication segment")
		i.GPInitSystem.ReplicationPort, _ = doWeHavePortBase(Config.INSTALL.FUTUREREFDIR, i.portFileName, "REPLICATION_PORT")
		if i.GPInitSystem.ReplicationPort == "" {
			Warnf("Didn't find REPLICATION_PORT in the file, setting it to default value: %d", defaultReplicatePort)
			i.GPInitSystem.ReplicationPort = strconv.Itoa(defaultReplicatePort)
		}
		i.GPInitSystem.ReplicationPort = strconv.Itoa(i.checkPortIsUsable(i.GPInitSystem.ReplicationPort))
	}
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
		// TODO: build multi
	}
}

// The contents of single node gpinitsystem
func (i *Installation) singleNodeGpInitSystem() []string{
	var primaryDir string
	for i := 1; i <= Config.INSTALL.TOTALSEGMENT; i++ {
		primaryDir = primaryDir + " " + Config.INSTALL.SEGMENTDATADIRECTORY
	}
	return []string{
		"ARRAY_NAME=" + i.GPInitSystem.ArrayName,
		"SEG_PREFIX=" + i.GPInitSystem.SegPrefix,
		"MASTER_HOSTNAME=" + i.GPInitSystem.MasterHostname,
		"MASTER_DIRECTORY=" + i.GPInitSystem.MasterDir,
		"PORT_BASE=" + i.GPInitSystem.SegmentPort,
		"MASTER_PORT=" + i.GPInitSystem.MasterPort,
		"DATABASE_NAME=" + i.GPInitSystem.DBName,
		"declare -a DATA_DIRECTORY=("+ primaryDir +")",
	}
}

// Check if the port is available or not
func (i *Installation) checkPortIsUsable(port string) int {
	pb, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
	}
	pb, err = isPortUsed(pb, Config.INSTALL.TOTALSEGMENT, i.WorkingHostFileLocation)
	if err != nil {
		fmt.Println(err)
	}
	return pb
}

func (i *Installation) executeGpInitSystem() {
	Infof("Executing the gpinitsystem to install the database")
	executeOsCommand(fmt.Sprintf("%s/bin/gpinitsystem", os.Getenv("GPHOME")), "-c", i.GpInitSystemConfigLocation, "-h", i.WorkingHostFileLocation , "-a")
}