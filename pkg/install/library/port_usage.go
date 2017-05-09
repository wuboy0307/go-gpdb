package library

import (
	"strconv"
	"net"
	log "../../core/logger"
	"os/exec"
	"strings"
	"../../core/methods"
	"../../core/arguments"
	"../objects"
)

// Check if the port is used, if yes then what is the next sequence program can use
func IsPortUsed(port int, iteration int) (int, error) {

	// Storing the base port
	BASE := port

	// Iterate through iteration to find how many port do we need
	for i:=1; i<=iteration; i++ {

		// Error out if the ports are not in the format needed
		_, err := strconv.ParseUint(strconv.Itoa(port), 10, 16)
		if err != nil { return 0, err }

		// Check if the port is available, if not find the next sequence
		ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
		if err != nil {
			log.Warn("PORT \"" + strconv.Itoa(port) + "\" is unavailable, finding the next sequence" )
			BASE =  BASE + iteration
			return IsPortUsed(BASE, iteration)
		}

		// Close listening to the port
		err = ln.Close()
		if err != nil {
			if err != nil { return 0, err }
		}

		// Iterate
		port =  port + i
	}

	// Return the collected information
	return BASE, nil

}

// Check if we have the last used ports
func DoWeHavePortBase(file string, name string, which_port string) (string, error) {

	portBaseFile := file + name
	return_code, err := methods.DoesFileOrDirExists(portBaseFile)
	if err != nil { return "", err }
	if return_code {
		log.Println("Found port file: \"" + portBaseFile + "\"")
		cmdOut, _ := exec.Command("grep", which_port, portBaseFile).Output()
		s := string(cmdOut)
		if strings.Contains(s, which_port) {
			return strings.TrimSpace(s), nil
		}
	} else {
		methods.CreateFile(portBaseFile)
		return "", nil
	}

	return "", nil
}


// Store the last used port
func StoreLastUsedPort() error {

	// Fully qualified filename
	FurtureRefFile := arguments.FutureRefDir + objects.PortBaseFileName
	log.Println("Storing the last used ports for this installation at: " + FurtureRefFile)

	// Delete the file if already exists
	err := methods.DeleteFile(FurtureRefFile)
	if err != nil {return nil}
	err = methods.CreateFile(FurtureRefFile)
	if err != nil {return nil}

	// Contents to write to the file
	var FutureContents []string
	FutureContents = append(FutureContents, "PORT_BASE: " + strconv.Itoa(objects.GpInitSystemConfig.PortBase + arguments.EnvYAML.Install.TotalSegments))
	FutureContents = append(FutureContents, "MASTER_PORT: " + strconv.Itoa(objects.GpInitSystemConfig.MasterPort + 1))

	// Write to the file
	err = methods.WriteFile(FurtureRefFile, FutureContents)
	if err != nil {return nil}

	return nil
}