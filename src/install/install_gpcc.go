package install

import (
	log "../../pkg/core/logger"
	"../../pkg/install/library"
	"../../pkg/core/arguments"
	"fmt"
	"errors"
)

// GPCC Installer
func InstalSingleNodeGPCC()  error {

	// Check if the provided GPDB version environment file exists
	env, err := library.PrevEnvFile("choose")
	if err != nil { return err }
	if env == "" {
		return errors.New("Couldn't find any environment file for the database version: " + arguments.RequestedInstallVersion + ", exiting...")
	}

	// Check if the binaries exists on the directory
	// if yes, Unzip the binaries if its file is zipped
	binary_file, err := UnzipBinary(arguments.RequestedCCInstallVersion)
	if err != nil { return err }
	fmt.Println(binary_file)

	// Check if the database is running

	// If exists then is there a GPCC already installed, if yes then ask for confirmation

	// If the request if to proceed, then uninstall the old GPCC installation.

	// Install the command center

	// Verify the command center is properly installed

	// Check the health of the database

	// Install the GPCC Web UI without WLM

	// Start the GPCC Web UI

	// Check if the ports are within the limit of unix port range

	// Store the last used port

	log.Println("Installation of GPCC software has been completed successfully")

	return nil
}