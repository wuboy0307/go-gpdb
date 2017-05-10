package install

import (
	"../../pkg/install/library"
	"../../pkg/install/objects"
	"../../pkg/core/arguments"
	"../../pkg/core/methods"
	"errors"
	log "../../pkg/core/logger"
	"strings"
	"strconv"
	"time"
	"os"
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
		log.Println("Using the environment file \"" + objects.EnvFileName + "\" for installing the GPCC")
	}

	// store the this database port and the GPHOME location
	err = library.ExtractEnvVariables(objects.EnvFileName)
	if err != nil { return err }

	// extract current date and time
	t := time.Now().Format("20060102150405")

	// If exists then is there a GPCC already installed, if yes then ask for confirmation
	if objects.GPCC_INSTANCE_NAME != "" {
		log.Warn("Found a instance of GPCC already installed on this environment, please confirm if the existing GPCC can be uninstalled")
		// Now ask for the confirmation
		confirm := methods.YesOrNoConfirmation()

		// What was the confirmation
		if confirm == "y" {  // yes, then uninstall the old GPCC installation
			err := library.UninstallGPCC(t, objects.EnvFileName)
			if err != nil { return err }
		} else { // no then exit
			log.Println("Cancelling the installation...")
			os.Exit(0)
		}
	}

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
	err = library.WasGPCCInstallationSucess()
	if err != nil { return err }

	// Whats is the next available port
	gpcc_port, err := library.DoWeHavePortBase(arguments.FutureRefDir, objects.GPCCPortBaseFileName, "GPCC_PORT")
	if err != nil { return err }
	if gpcc_port != "" {
		objects.GPCC_PORT = string(gpcc_port)[11:]
	} else {
		log.Warn("Didn't find GPCC_PORT in the file, setting it to default value: " + strconv.Itoa(objects.GPCC_PORT_BASE))
		objects.GPCC_PORT = strconv.Itoa(objects.GPCC_PORT_BASE)
	}

	// Is the port used
	ccp, err := strconv.Atoi(objects.GPCC_PORT)
	if err != nil { return err }
	ccp, err = library.IsPortUsed(ccp, 1)
	if err != nil { return err }
	log.Println("Setting the GPCC_PORT has: "+ strconv.Itoa(ccp))
	objects.GPCC_PORT = strconv.Itoa(ccp)

	// Install the GPCC Web UI without WLM
	cc_name := "gpcc_" + arguments.RequestedInstallVersion + "_" + arguments.RequestedCCInstallVersion + "_" + t
	cc_name = strings.Replace(cc_name, ".", "", -1)
	if strings.HasPrefix(arguments.RequestedCCInstallVersion, "1") {
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(objects.ThisDBMasterPort), objects.GPCC_PORT, "n", "n", "n", "n", "EOF"}
		library.InstallGPCCUI(script_option, objects.BinaryInstallLocation)
	} else if strings.HasPrefix(arguments.RequestedCCInstallVersion, "2") {
		objects.InstallWLM = true
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(objects.ThisDBMasterPort), "n", objects.GPCC_PORT, "n", "n", "n", "n", "EOF"}
		library.InstallGPCCUI(script_option, objects.BinaryInstallLocation)
	} else {
		objects.InstallWLM = true
		var script_option = []string{cc_name, cc_name, "n", strconv.Itoa(objects.ThisDBMasterPort), "n", objects.GPCC_PORT, "n", "n", "EOF"}
		library.InstallGPCCUI(script_option, objects.BinaryInstallLocation)
	}

	// Start the GPCC Web UI
	err = library.StartGPCC(cc_name, objects.BinaryInstallLocation)
	if err != nil { return err }

	// If success update the environment file
	err = library.UpdateEnvFile(cc_name)

	// Check if the port is not greater than 63000, since unix limit is 64000
	port, _ := strconv.Atoi(objects.GPCC_PORT)
	if  port > 63000 {
		log.Warn("PORT has execeeded the unix port limit, setting it to default.")
		objects.GPCC_PORT = strconv.Itoa(objects.GPCC_PORT_BASE)
	}

	// Store the last used port
	err = library.StoreLastUsedGPCCPort()
	if err != nil { return err }

	//// Uninstall the WLM
	if !methods.IsValueEmpty(objects.WLMInstallDir) {
		err := library.UninstallWLM(t)
		if err != nil { return err }
	}

	// Install the WLM
	err = library.InstallWLM()
	if err != nil { return err }

	// Start the new browser for use

	// when starting , start gpcc instance & WLM

	log.Println("Installation of GPCC software has been completed successfully")

	return nil
}