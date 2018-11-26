package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Create environment file of this installation
func (i *Installation) createEnvFile() {

	// Environment file fully qualified path
	i.EnvFile = Config.INSTALL.ENVDIR + "env_" + cmdOptions.Version + "_" + i.Timestamp
	Infof("Creating environment file for this installation at: " + i.EnvFile)

	// Create the file
	createFile(i.EnvFile)
	writeFile(i.EnvFile, []string{
		"export GPHOME="+ i.BinaryInstallationLocation,
		"export PYTHONPATH=$GPHOME/lib/python",
		"export PYTHONHOME=$GPHOME/ext/python",
		"export PATH=$GPHOME/bin:$PYTHONHOME/bin:$PATH",
		"export LD_LIBRARY_PATH=$GPHOME/lib:$PYTHONHOME/lib:$LD_LIBRARY_PATH",
		"export OPENSSL_CONF=$GPHOME/etc/openssl.cnf",
		"export MASTER_DATA_DIRECTORY="+ i.GPInitSystem.MasterDir + "/" + i.GPInitSystem.ArrayName + "-1",
		"export PGPORT=" + i.GPInitSystem.MasterPort,
		"export PGDATABASE=" + i.GPInitSystem.DBName,
	})
}

// Update Environment file
func (i *Installation) updateEnvFile() error {

	Infof("Updating the environment file \"%s\" with the GPCC environment", i.EnvFile)

	// Append to file
	appendFile(i.EnvFile, []string{
		"export GPPERFMONHOME=" + "sdafg",
		"export PATH=$GPPERFMONHOME/bin:$PATH",
		"export LD_LIBRARY_PATH=$GPPERFMONHOME/lib:$LD_LIBRARY_PATH",
		"export GPCC_INSTANCE_NAME=" + i.GPCC.InstanceName,
		"export GPCCPORT=" + i.GPCC.InstancePort,
	})

	return nil
}

// Read file and find the content that we are interested
func readFileAndGatherInformation(file string) string {

	Debugf("Reading and gathering information for the file: %s", file)

	// Obtain the below information from env file
	var output string

	// From the file find the text detection information
	content := readFile(file)

	// DB PORT
	c := contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", "export PGPORT="), []string{"FS", "="})
	output = output + "|" + removeBlanks(c.String())

	// Is DB running
	if isDbHealthy(file, "") {
		output = output + "| RUNNING"
	} else {
		output = output + "| STOPPED"
	}

	// GPCC Instance
	c = contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", "export GPCC_INSTANCE_NAME="), []string{"FS", "="})
	output = output + "|" + removeBlanks(c.String())

	// GPCC URL
	c = contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", "export GPCCPORT="), []string{"FS", "="})
	if !IsValueEmpty(c.String()) {
		output = output + "|" + fmt.Sprintf("http://%s:%s", GetLocalIP(), removeBlanks(c.String()))
	}

	return output
}

// List all the installed enviornment files.
func installedEnvFiles(search, confirmation string, ignoreErr bool) string {

	Debugf("Searching for installed env file using the search string: %s", search)

	var output = []string{`Index | Environment File | Master Port | Status | GPCC Instance Name | GPCC Instance URL`,
		`------|-----------------------|-----------------|------------------|----------------------------------------|------------------------------------------`,
	}

	// Search for environment files for that version
	allEnv, _ := FilterDirsGlob(Config.INSTALL.ENVDIR, search)
	if (len(allEnv) == 0 && !IsValueEmpty(cmdOptions.Version)) && !ignoreErr {
		Fatalf("No installation found for the version: %s, try downloading and installing it", cmdOptions.Version)
	} else if len(allEnv) == 0 && IsValueEmpty(cmdOptions.Version) {
		Fatalf("No installation found, try downloading and installing it")
	}

	// All the environments
	for k, v := range allEnv {
		envName := strings.Replace(v, Config.INSTALL.ENVDIR, "", -1)
		output = append(output, fmt.Sprintf("%s|%s%s", strconv.Itoa(k+1), envName, readFileAndGatherInformation(v)))
	}

	// Found matching environment file of this installation, now ask for confirmation
	if len(allEnv) > 0 && confirmation == "confirm" {

		printOnScreen("Here is a list of env installed, confirm if you want to continue", output)

		// Now ask for the confirmation
		confirm := YesOrNoConfirmation()

		// What was the confirmation
		if confirm == "y" { // yes
			Infof("Continuing with the installation of version: %s", cmdOptions.Version)
			return allEnv[0]
		} else { // no
			Infof("Cancelling the installation...")
			os.Exit(0)
		}

	} else if len(allEnv) == 1 && confirmation == "choose" { // if there is only one , then there is no choose just provide the only one
		startDBifNotStarted(allEnv[0])
		return allEnv[0]

	} else if (len(allEnv) > 0 && confirmation == "choose") || confirmation == "list&choose" {

		printOnScreen(fmt.Sprintf("Found %d installation, choose from the list", len(allEnv)), output)

		// What is users choice
		choice := PromptChoice(len(allEnv))

		// return the enviornment file to the main function
		choosenEnv := allEnv[choice-1]
		startDBifNotStarted(choosenEnv)

		return choosenEnv
	}

	return ""
}

// List all the environment installed on this box
func env() {
	var envFile string
	if cmdOptions.Version == "" { // No version provided, show everything
		Infof("Listing all the environment installed")
		envFile = installedEnvFiles("*", "list&choose", false)
	} else { // Version given, search for env file
		Infof("Listing all the environment installed with version: %s", cmdOptions.Version)
		envFile = installedEnvFiles("*" + cmdOptions.Version + "*", "choose", false)
	}

	displayEnvFileToSource(envFile)
}

// Display the env content on the screen
func displayEnvFileToSource(file string) {
	printOnScreen("Source the environment file to set the environment", []string{"source " + file})
}