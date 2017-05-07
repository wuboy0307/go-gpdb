package install

import (
	"../../pkg/install/library"
	"../../pkg/install/objects"
	"../../pkg/core/arguments"
	"errors"
	log "../../pkg/core/logger"
)

// GPCC Installer
func InstalSingleNodeGPCC()  error {

	// Check if the provided GPDB version environment file exists
	env, err := library.PrevEnvFile("choose")
	if err != nil { return err }
	if env == "" {
		return errors.New("Couldn't find any environment file for the database version: " + arguments.RequestedInstallVersion + ", exiting...")
	} else {
		objects.EnvFileName = arguments.EnvFileDir + env
	}

	// store the this database port and the GPHOME location
	err = library.ExtractPortAndGPHOME(objects.EnvFileName)
	if err != nil { return err }

	// Check if the database is running, if not then start the database
	err = library.StartDBifNotStarted()
	if err != nil { return err }

	// Check if the binaries exists on the directory
	// if yes, Unzip the binaries if its file is zipped
	binary_file, err := UnzipBinary(arguments.RequestedCCInstallVersion)
	if err != nil { return err }

	// Extract the binaries.
	objects.BinaryInstallLocation = "/usr/local/greenplum-cc-web-dbv-" + arguments.RequestedInstallVersion + "-ccv-" + arguments.RequestedCCInstallVersion
	var script_option = []string{"yes", objects.BinaryInstallLocation, "yes", "yes"}
	err = library.ExecuteBinaries(binary_file, objects.InstallGPDBBashFileName, script_option)
	if err != nil { return err }

	// If exists then is there a GPCC already installed, if yes then ask for confirmation

	// If the request if to proceed, then uninstall the old GPCC installation.

	// Install the command center database (GPPERFMON)
	err = library.InstallGpperfmon()
	if err != nil { return err }

	// Restart the database to make the changes to take into effect
	log.Println("Restarting the database so that the command center agents can be started")
	err = library.StopDB()
	if err != nil { return err }
	err = library.StartDB()
	if err != nil { return err }

	// Verify the command center is properly installed

	// Check the health of the database

	// Install the GPCC Web UI without WLM

	// Start the GPCC Web UI

	// Check if the ports are within the limit of unix port range

	// Store the last used port

	log.Println("Installation of GPCC software has been completed successfully")

	return nil
}