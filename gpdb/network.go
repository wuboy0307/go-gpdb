package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// validate port
func (i *Installation) validatePort(searchString string, defaultPort int) string {
	Infof("Obtaining ports to be set for %s", searchString)
	p, _ := doWeHavePortBase(Config.INSTALL.FUTUREREFDIR, i.PortFileName, searchString)
	if p == "" {
		Warnf("Didn't find %s in the file, setting it to default value: %d", searchString, defaultPort)
		p = strconv.Itoa(defaultPort)
	}
	p = strconv.Itoa(i.checkPortIsUsable(p))
	return p
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

	if cmdOptions.Product == "gpdb" {
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
		portFile := Config.INSTALL.FUTUREREFDIR + i.PortFileName
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
	} else { // gpcc
		// Check if the port is not greater than 63000, since unix limit is 64000
		if outOfRangePort(i.GPCC.InstancePort) || outOfRangePort(i.GPCC.WebSocketPort)  {
			Warnf("PORT has exceeded the unix port limit, setting it to default.")
			i.GPCC.InstancePort = strconv.Itoa(defaultGpccPort)
			i.GPCC.WebSocketPort = strconv.Itoa(defaultWebSocket)
		}

		// Fully qualified filename
		portFile := Config.INSTALL.FUTUREREFDIR + i.PortFileName
		Infof("Storing the last used ports for this installation at: %s", portFile)

		// Delete the file if already exists
		deleteFile(portFile)
		createFile(portFile)
		writeFile(portFile, []string{
			"GPCC_PORT: " + strconv.Itoa(strToInt(i.GPCC.InstancePort) + 1),
			"WEBSOCKET_PORT: " + strconv.Itoa(strToInt(i.GPCC.WebSocketPort) + 1),
		})
	}
}

// Parse the IP and see if they are in the range
var CoreIP = "192.0.0.0"
func checkIP(found_ip string) bool {
	var ip = net.ParseIP(CoreIP)

	// Parse the IP
	obtain_ip := net.ParseIP(found_ip)

	// Get the local IP by comparing it.
	if bytes.Compare(obtain_ip, ip) >= 0 {
		return true
	}

	return false
}

// Get the IP address
func GetLocalIP() (string){

	Infof("Getting local IP address")

	// Get Interface address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	// Extract the IP's and see if its in the range.
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if checkIP(ipnet.IP.String()) {
					Infof("Found Local IP address: %s", ipnet.IP.String())
					return ipnet.IP.String()
				}
			}
		}
	}

	// If nothing found then produce a warning message
	Warnf("Cannot find local IP address in the range \"%s\", ignoring the IP....", CoreIP)
	return ""
}