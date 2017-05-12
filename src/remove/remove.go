package remove

import (
	log "../../pkg/core/logger"
	"../../pkg/core/methods"
	"../../pkg/install/library"
	"../../pkg/install/objects"
	"../../pkg/core/arguments"
	"strings"
)

// Run gpdeletesystem
func deleteGPDBEnvUsingGpDeleteSystem() error {

	log.Println("Calling gpdeletesystem to remove the environment")

	// Script arguments
	var deleteSystemArgs []string
	deleteSystemArgs = append(deleteSystemArgs, "source " + objects.EnvFileName)
	deleteSystemArgs = append(deleteSystemArgs, "gpdeletesystem -d $MASTER_DATA_DIRECTORY -f << EOF")
	deleteSystemArgs = append(deleteSystemArgs, "y")
	deleteSystemArgs = append(deleteSystemArgs, "y")
	deleteSystemArgs = append(deleteSystemArgs, "EOF")

	// Write it to the file.
	file := arguments.TempDir + "run_deletesystem.sh"
	err := library.ExecuteBash(file, deleteSystemArgs)
	if err != nil { return err }

	return nil
}


// Cleaning up all the environments files etc
func deleteGPDBEnvUsingManualMethod(version string, timestamp string) error {

	log.Println("Running cleanup script to clean the environment listing")

	// Script arguments
	var removeAllargs []string
	removeAllargs = append(removeAllargs, "/bin/sh " + arguments.UninstallDir + "uninstall_" + version + "_" + timestamp)

	// Write it to the file.
	file := arguments.TempDir + "run_cleanup_script.sh"
	err := library.ExecuteBash(file, removeAllargs)
	if err != nil { return err }

	return nil
}

// Main Remove method
func Remove(version string) error {

	log.Println("Starting program to uninstall the version: " + version)

	// store the variable
	arguments.RequestedInstallVersion = version

	// Check if the envfile for that version exists
	chosenEnvFile, err := library.PrevEnvFile("choose")
	if err != nil { return err }

	// If we receive none, then display the error to user
	var timestamp string
	if methods.IsValueEmpty(chosenEnvFile) {
		log.Fatal("Cannot find any environment with the version: " + version )
	} else { // Else store the value
		objects.EnvFileName = arguments.EnvFileDir + chosenEnvFile
		timestamp = strings.Split(chosenEnvFile, "_")[2]
	}

	log.Println("The choosen enviornment file to remove is: " + objects.EnvFileName)

	// store the this database port and the GPHOME location
	err = library.ExtractEnvVariables(objects.EnvFileName)
	if err != nil { return err }

	// Check if the database is running, if not then start the database
	err = library.StartDBifNotStarted()
	if err != nil { return err }

	// Cleanup GPCC if installed.
	if !methods.IsValueEmpty(objects.GPPERFMONHOME) {
		err = library.UninstallGPCC(timestamp, objects.EnvFileName)
		if err != nil { return err }
	}

	// run delete command
	err = deleteGPDBEnvUsingGpDeleteSystem()
	if err != nil {
		log.Warn("gpdeletesystem failed, running manual cleanup")
		err = deleteGPDBEnvUsingManualMethod(version, timestamp)
		if err != nil { return err }
	} else {
		log.Println("Deleted GPDB via gpdeletesystem was a success.")
		err = deleteGPDBEnvUsingManualMethod(version, timestamp)
		if err != nil { return err }
	}

	log.Println("Uninstallation of environment \"" + chosenEnvFile + "\" was a success")
	log.Println("exiting ....")

	return nil
}
