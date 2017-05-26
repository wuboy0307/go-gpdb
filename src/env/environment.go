package env

import (
	log "../../pkg/core/logger"
	"../../pkg/install/library"
	"../../pkg/install/objects"
	"../../pkg/core/arguments"
	"../../pkg/core/methods"
	"errors"
)

// Set the environment
func SettheChoosenEnv(chosenEnvFile string, version string) error {
	// If we receive none, then display the error to user
	if methods.IsValueEmpty(chosenEnvFile) {
		log.Fatal("Cannot find any environment with the version: " + version )
	} else { // Else store the value
		objects.EnvFileName = arguments.EnvFileDir + chosenEnvFile
	}

	log.Println("The choosen enviornment file is: " + objects.EnvFileName)

	// store the this database port and the GPHOME location
	err := library.ExtractEnvVariables(objects.EnvFileName)
	if err != nil { return err }

	// Check if the database is running, if not then start the database
	err = library.StartDBifNotStarted()
	if err != nil { return err }

	log.Println("Environment has been setup and ready to use")

	// Open terminal after source the environment file
	err = library.SetVersionEnv(objects.EnvFileName)
	if err != nil { return err }

	log.Println("Environment setup is complete")

	return nil
}


// Function to setup environment
func Environment(version string) error {

	// If no version is provided then the user is requesting to list all the environment installed
	if version == "" {
		log.Println("listing all the environment version installed on this cluster")
		chosenEnvFile, err := library.PrevEnvFile("listandchoose")
		if err != nil { return err }

		// set the chosen environment
		if methods.IsValueEmpty(chosenEnvFile) {
			return errors.New("There is no any installation of GPDB, please install the product to list the environment here")
		} else {
			err = SettheChoosenEnv(chosenEnvFile, version)
			if err != nil { return err }
		}


	} else { // he is checking for a specific version

		log.Println("listing all the environment that has been installed with version: " + version)

		// store the variable
		arguments.RequestedInstallVersion = version

		// Get the env files that we know about
		chosenEnvFile, err := library.PrevEnvFile("choose")
		if err != nil { return err }

		// set the chosen environment
		if methods.IsValueEmpty(chosenEnvFile) {
			return errors.New("There is no any installation of GPDB, please install the product to list the environment here")
		} else {
			err = SettheChoosenEnv(chosenEnvFile, version)
			if err != nil { return err }
		}
	}

	log.Println("Exiting ..... ")
	return nil
}
