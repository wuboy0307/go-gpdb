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

// Generate the port for master / segments / mirror & replication
func (i *Installation) generatePortRange() {

	// Check if we have the last used port base file and its usable
	i.GPInitSystem.SegmentPort = i.validatePort(Config.INSTALL.FUTUREREFDIR, "PRIMARY_PORT", defaultPrimaryPort)  // segment
	i.GPInitSystem.MasterPort = i.validatePort(Config.INSTALL.FUTUREREFDIR, "MASTER_PORT", defaultMasterPort) // master

	// If its a multi installation we will need the mirror / repliocateion port as well & usable
	if i.SingleORMulti == "multi" {
		i.GPInitSystem.MirrorPort = i.validatePort(Config.INSTALL.FUTUREREFDIR, "MIRROR_PORT", defaultMirrorPort) // mirror
		i.GPInitSystem.ReplicationPort = i.validatePort(Config.INSTALL.FUTUREREFDIR, "REPLICATION_PORT", defaultReplicatePort) // replication
	}
}

// validate port
func (i *Installation) validatePort(dir, searchString string, defaultPort int) string {
	Infof("Obtaining ports to be set for %s", searchString)
	p, _ := doWeHavePortBase(dir, i.portFileName, searchString)
	if p == "" {
		Warnf("Didn't find %s in the file, setting it to default value: %d", searchString, defaultPort)
		p = strconv.Itoa(defaultPort)
	}
	p = strconv.Itoa(i.checkPortIsUsable(i.GPInitSystem.SegmentPort))
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