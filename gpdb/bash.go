package main

import (
	"fmt"
	"os/exec"
)

// Build executable shell script
func buildExecutableBashScript(filename string, executableFilename string, args []string) error {

	Debugf("Starting program to build executable script to install the binaries")

	// Create the file
	createFile(filename)

	// Create arguments on what needs to be written to the file
	Infof("Generating installer arguments to be passed to the file: %s", filename)
	executableLine := fmt.Sprintf("/bin/sh %s &>/dev/null << EOF", executableFilename)

	// Load the contents to the file
	var passArgs []string
	passArgs = append(passArgs, executableLine)
	for _, v := range args {
		passArgs = append(passArgs, v)
	}
	passArgs = append(passArgs, "EOF")

	// Write to the file
	writeFile(filename, passArgs)

	return nil
}

// Execute shell script when called
func executeBinaries(binaryFile string, bashfilename string, scriptOptions []string) error {

	// Build a quick shell script to install binaries
	// Filename name
	filename := Config.CORE.TEMPDIR + bashfilename

	// Cleanup the file if already exists (ignore error if any)
	deleteFile(filename)

	// Create the shell script
	err := buildExecutableBashScript(filename, binaryFile, scriptOptions)
	if err != nil {
		return err
	}

	// Execute the installer script
	Infof("Executing the bash script: %s", filename)
	_, err = exec.Command("/bin/sh", filename).Output()
	if err != nil {
		return err
	}

	// Cleanup the tempFile
	Infof("Cleaning up the file \"%s\", if found", filename)
	deleteFile(filename)

	return nil
}
