package library

import (
	log "../../core/logger"
	"../../core/arguments"
	"../../core/methods"
	"../objects"
	"strconv"
	"io/ioutil"
	"strings"
	"fmt"
	"bufio"
	"os"
	"os/exec"
)

// Create environment file of this installation
func CreateEnvFile(t string) error {

	// Environment file fully qualified path
	objects.EnvFileName = arguments.EnvFileDir + "env_" + arguments.RequestedInstallVersion + "_" + t
	log.Println("Creating environment file for this installation at: " + objects.EnvFileName)

	// Create the file
	err := methods.CreateFile(objects.EnvFileName)
	if err != nil { return err }

	// Build arguments to write
	var EnvFileContents []string
	EnvFileContents = append(EnvFileContents, "source " + objects.BinaryInstallLocation + "/greenplum_path.sh")
	EnvFileContents = append(EnvFileContents, "export MASTER_DATA_DIRECTORY=" + objects.GpInitSystemConfig.MasterDir + "/" + objects.GpInitSystemConfig.ArrayName + "-1")
	EnvFileContents = append(EnvFileContents, "export PGPORT=" + strconv.Itoa(objects.GpInitSystemConfig.MasterPort))
	EnvFileContents = append(EnvFileContents, "export PGDATABASE=" + objects.GpInitSystemConfig.DBName)

	// Write to EnvFile
	err = methods.WriteFile(objects.EnvFileName, EnvFileContents)
	if err != nil { return err }

	return nil
}

// Check if there is any previous installation of the same version
func PrevEnvFile() error {

	log.Println("Checking if there is previous installation for the version: " + arguments.RequestedInstallVersion)
	var MatchingFilesInDir []string
	allfiles, err := ioutil.ReadDir(arguments.EnvFileDir)
	if err != nil { return err }
	for _, file := range allfiles {

		if strings.Contains(file.Name(), arguments.RequestedInstallVersion) {
			MatchingFilesInDir = append(MatchingFilesInDir, file.Name())
		}

	}

	// Found matching environment file of this installation, now ask for confirmation
	if len(MatchingFilesInDir) > 1 {

		// Show all the environment files
		log.Warn("Found matching environment file for the version: " + arguments.RequestedInstallVersion)
		log.Println("Below are the list of environment file of the version: " + arguments.RequestedInstallVersion + "\n")
		for _, k := range MatchingFilesInDir {
			fmt.Printf("%v \n", k)
		}
		fmt.Println()

		// Ask for confirmation
		confirm := YesOrNoConfirmation()

		// What was the confirmation
		if confirm == "y" {  // yes
			log.Println("Continuing with the installtion of version: " + arguments.RequestedInstallVersion)
		} else { // no
			log.Println("Cancelling the installation...")
			os.Exit(0)
		}
	}
	return nil
}

// Prompt for confirmation
func YesOrNoConfirmation() string {

	// Start the new scanner to get the user input
	fmt.Print("You can use \"gpdb env -v <version>\" to set the env, do you wish to continue (Yy/Nn)?: ")
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {

		// The choice entered
		choice_entered := input.Text()

		// If its a valid value move on
		if arguments.YesOrNo[strings.ToLower(choice_entered)] == "y" {  // Is it Yes
			return choice_entered
		} else if arguments.YesOrNo[strings.ToLower(choice_entered)] == "n" { // Is it No
			return choice_entered
		} else { // Invalid choice, ask to re-enter
			fmt.Println("Invalid Choice: Please enter Yy/Nn, try again.")
			return YesOrNoConfirmation()
		}
	}

	return ""
}

// Set Environment of the shell
func SetVersionEnv(filename string) error {

	log.Println("Attempting to open a terminal, after setting the environment of this installation.")

	// User Home
	usersHomeDir := os.Getenv("HOME")

	// Create a temp file to execute
	executeFile := arguments.TempDir + "openterminal.sh"
	_ = methods.DeleteFile(executeFile)
	_ = methods.CreateFile(executeFile)

	// The command
	var cmd []string
	cmdString := "gnome-terminal --working-directory=\"" + usersHomeDir + "\" --tab -e 'bash -c \"echo \\\"Sourcing Envionment file: "+ filename + "\\\"; source "+ filename +"; exec bash\"'"
	cmd = append(cmd, cmdString)

	// Write to the file
	_ = methods.WriteFile(executeFile, cmd)
	_, err := exec.Command("/bin/sh", executeFile).Output()
	if err != nil { return nil }

	// Cleanup the file file.
	_ = methods.DeleteFile(executeFile)

	return nil
}