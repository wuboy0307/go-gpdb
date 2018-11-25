package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Check if the port 22 is reachable, should return back in 5 seconds
func checkHostReachability(address string, errorType bool) bool {
	Debugf("Checking the connectivity for the address: %s", address)
	conn, err := net.DialTimeout("tcp", address, time.Second * 5)
	if err != nil {
		if errorType { // If we care about the error
			Errorf("Could not reach the host: %s", address)
		} else { // Here we are looking for free port, so place it on debug
			Debugf("The host and the port \"%s\" is not used and can be used", address)
		}
		return false
	}
	defer conn.Close()
	return true
}

// Check if the port is used, if yes then what is the next sequence program can use
func isPortUsed(port int, iteration int, hostfileLoc string) (int, error) {

	Debugf("Checking for the port usage %d", port)

	// Storing the base port
	BASE := port

	// Iterate through iteration to find how many port do we need
	for i := 1; i <= iteration; i++ {

		// Error out if the ports are not in the format needed
		_, err := strconv.ParseUint(strconv.Itoa(port), 10, 16)
		if err != nil {
			return 0, fmt.Errorf("err in parsing the integer, err: %v", err)
		}

		// Check if the port is available, if not find the next sequence
		for _, host := range strings.Split(string(readFile(hostfileLoc)), "\n"){
			if host != "" && checkHostReachability(fmt.Sprintf("%s:%d", host, port), false) {
				Warnf("PORT \"%d\" is unavailable, finding the next sequence", port)
				BASE = BASE + iteration
				return isPortUsed(BASE, iteration, hostfileLoc)
			}
		}

		// Iterate
		port = port + i
	}

	// Return the collected information
	return BASE, nil

}

// Check if we have the last used ports
func doWeHavePortBase(file string, name string, whichPort string) (string, error) {

	// the port base file
	portBaseFile := file + name
	Debugf("Checking for the port for \"%s\" from the file: %s", whichPort, portBaseFile)
	returnCode, err := doesFileOrDirExists(portBaseFile)
	if err != nil {
		Fatalf("Error when checking for port base, err: %v", err)
	}

	// Extract the port if found, else create a file.
	if returnCode {
		Debugf("Found port file: %s", portBaseFile)
		port := contentExtractor(readFile(portBaseFile), fmt.Sprintf("/^%s/ {print $2}", whichPort), []string{"FS", ":"})
		return removeBlanks(port.String()), nil
	} else {
		createFile(portBaseFile)
	}

	return "", nil
}

// Save the last used port number
func (i *Installation) savePort() {

	// Check if the port is not greater than 63000, since unix limit is 64000
	if outOfRangePort(i.GPInitSystem.SegmentPort) || outOfRangePort(i.GPInitSystem.MasterPort) || outOfRangePort(i.GPInitSystem.ReplicationPort) || outOfRangePort(i.GPInitSystem.MirrorPort) || outOfRangePort(i.GPInitSystem.MirrorReplicationPort) {
		Warnf("PORT has exceeded the unix port limit, setting it to default.")
		i.GPInitSystem.SegmentPort = strconv.Itoa(defaultPrimaryPort)
		i.GPInitSystem.MasterPort = strconv.Itoa(defaultMasterPort)
		i.GPInitSystem.ReplicationPort = strconv.Itoa(defaultReplicatePort)
		i.GPInitSystem.MirrorPort = strconv.Itoa(defaultMirrorPort)
		i.GPInitSystem.MirrorReplicationPort = strconv.Itoa(defaultMirrorReplicatePort)
	}

	// Fully qualified filename
	portFile := Config.INSTALL.FUTUREREFDIR + i.portFileName
	Infof("Storing the last used ports for this installation at: %s", portFile)

	// Delete the file if already exists
	deleteFile(portFile)
	createFile(portFile)
	writeFile(portFile, []string{
		"PRIMARY_PORT: " + strconv.Itoa(strToInt(i.GPInitSystem.SegmentPort) + Config.INSTALL.TOTALSEGMENT),
		"MASTER_PORT: " + strconv.Itoa(strToInt(i.GPInitSystem.MasterPort) + 1),
		"REPLICATION_PORT: " + strconv.Itoa(strToInt(i.GPInitSystem.ReplicationPort) + Config.INSTALL.TOTALSEGMENT),
		"MIRROR_PORT: " + strconv.Itoa(strToInt(i.GPInitSystem.MirrorPort) + Config.INSTALL.TOTALSEGMENT),
		"MIRROR_REPLICATION_PORT: " + strconv.Itoa(strToInt(i.GPInitSystem.MirrorReplicationPort) + Config.INSTALL.TOTALSEGMENT),
	})
}