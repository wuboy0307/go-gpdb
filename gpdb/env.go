package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

// Read file and find the content that we are interested
func readFileAndGatherInformation(file string) string {

	// Obtain the below information from env file
	var textDetection = []string{`export PGPORT=`, `export GPCC_INSTANCE_NAME=`, `export GPCCPORT=`}
	var output string

	// Open the file
	f, err := os.Open(file)
	if err != nil {
		Fatalf("Error in opening the file, err: %v", err)
	}
	defer f.Close()

	// From the file find the text detection information
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _, text := range textDetection {
			if strings.HasPrefix(scanner.Text(), text) {
				output = output + "|" + strings.Replace(scanner.Text(), text, "", -1)
				if strings.HasPrefix(scanner.Text(),"export PGPORT=") {
					output = output + "| STOPPED"
				}
			}
		}
	}

	// If any error throw on the screen
	if err := scanner.Err(); err != nil {
		Fatalf("Error in scanning the file, err: %v", err)
	}

	return output
}

// List all the installed enviornment files.
func installedEnvFiles(search, confirmation string) string {

	var output = []string{`Index | Environment File | Master Port | Status | GPCC Instance Name | GPCC Instance URL`,
		`------|-----------------------|-----------------|------------------|----------------------------------------|------------------------------------------`,
	}

	// Search for environment files for that version
	allEnv, _ := FilterDirsGlob(Config.INSTALL.ENVDIR, search)
	if len(allEnv) == 0 && !IsValueEmpty(cmdOptions.Version) {
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

		return allEnv[0]

	} else if (len(allEnv) > 0 && confirmation == "choose") || confirmation == "list&choose" {

		printOnScreen(fmt.Sprintf("Found %d installation, choose from the list", len(allEnv)), output)

		// What is users choice
		choice := PromptChoice(len(allEnv))

		// return the enviornment file to the main function
		choosenEnv := allEnv[choice-1]

		return choosenEnv
	}

	return ""
}

func envListing() {

	var envFile string
	if cmdOptions.Version == "" { // No version provided, show everything
		Infof("Listing all the environment installed")
		envFile = installedEnvFiles("*", "list&choose")
	} else { // Version given, search for env file
		Infof("Listing all the environment installed with version: %s", cmdOptions.Version)
		envFile = installedEnvFiles("*" + cmdOptions.Version + "*", "choose")
	}

	printOnScreen("Source the environment file to set the environment", []string{"source " + envFile})
}
