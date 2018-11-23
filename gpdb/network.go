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

	// Storing the base port
	BASE := port

	// Iterate through iteration to find how many port do we need
	for i := 1; i <= iteration; i++ {

		// Error out if the ports are not in the format needed
		_, err := strconv.ParseUint(strconv.Itoa(port), 10, 16)
		if err != nil {
			return 0, err
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
	returnCode, err := doesFileOrDirExists(portBaseFile)
	if err != nil {
		Fatalf("Error when checking for port base, err: %v", err)
	}

	// Extract the port if found, else create a file.
	if returnCode {
		Debugf("Found port file: %s", portBaseFile)
		port := contentExtractor(readFile(portBaseFile), fmt.Sprintf("/%s/ {print $2}", whichPort), []string{"FS", "="})
		return removeBlanks(port.String()), nil
	} else {
		createFile(portBaseFile)
	}

	return "", nil
}

// Store the last used port
func storeLastUsedPort(filename string, masterPort, segmentPort int) error {

	// Fully qualified filename
	FurtureRefFile := Config.INSTALL.FUTUREREFDIR + filename
	Infof("Storing the last used ports for this installation at: %s", FurtureRefFile)

	// Delete the file if already exists
	deleteFile(FurtureRefFile)
	createFile(FurtureRefFile)

	// Contents to write to the file
	var FutureContents []string
	FutureContents = append(FutureContents, "PORT_BASE: "+strconv.Itoa(segmentPort+Config.INSTALL.TOTALSEGMENT))
	FutureContents = append(FutureContents, "MASTER_PORT: "+strconv.Itoa(masterPort+1))

	// Write to the file
	writeFile(FurtureRefFile, FutureContents)

	return nil
}
