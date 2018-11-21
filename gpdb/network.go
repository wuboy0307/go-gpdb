package main

import (
	"net"
	"time"
	"strconv"
	"os/exec"
	"strings"
)

// Check if the port 22 is reachable, should return back in 5 seconds
func checkHostReachability(address string) bool {
	_, err := net.DialTimeout("tcp", address, time.Second*5)
	if err != nil {
		Errorf("Could not reach the host: %s", address)
		return false
	}
	return true
}

// Check if the port is used, if yes then what is the next sequence program can use
func isPortUsed(port int, iteration int) (int, error) {

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
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			Warnf("PORT \"%s\" is unavailable, finding the next sequence", strconv.Itoa(port))
			BASE = BASE + iteration
			return isPortUsed(BASE, iteration)
		}

		// Close listening to the port
		err = ln.Close()
		if err != nil {
			if err != nil {
				return 0, err
			}
		}

		// Iterate
		port = port + i
	}

	// Return the collected information
	return BASE, nil

}

// Check if we have the last used ports
func doWeHavePortBase(file string, name string, which_port string) (string, error) {

	portBaseFile := file + name
	returnCode, err := doesFileOrDirExists(portBaseFile)
	if err != nil {
		return "", err
	}

	if returnCode {
		Infof("Found port file: %s", portBaseFile)
		cmdOut, _ := exec.Command("grep", which_port, portBaseFile).Output()
		s := string(cmdOut)
		if strings.Contains(s, which_port) {
			return strings.TrimSpace(s), nil
		}
	} else {
		createFile(portBaseFile)
		return "", nil
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
