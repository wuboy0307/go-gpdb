package main

import (
	"fmt"
	"os"
	"time"
)

func (i *Installation) preGPCCChecks() {

	Infof("Running the pre checks to install GPCC version \"%s\" on the GPDB Version \"%s\"", cmdOptions.CCVersion, cmdOptions.Version)
	// Check if there is already a version of GPDB installed
	i.EnvFile = installedEnvFiles(fmt.Sprintf("*%s*", cmdOptions.Version), "choose", false)

	// Extract the environment information
	i.extractEnvValues()

	// Extract current date and time
	i.Timestamp = time.Now().Format("20060102150405")

	// Check if this env have GPCC Installed
	i.doesThisEnvHasGPCCInstalled()

	// Check if the database is running, if not then start the database
	startDBifNotStarted(i.EnvFile)
}

// Extract environment values from the env file
func (i *Installation) extractEnvValues() {
	Infof("Extracting the environment information from the file: %s", i.EnvFile)
	content := readFile(i.EnvFile)
	c := contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", "export PGPORT="), []string{"FS", "="})
	i.GPInitSystem.MasterPort = removeBlanks(c.String())
	c = contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", "export MASTER_DATA_DIRECTORY="), []string{"FS", "="})
	i.GPInitSystem.MasterDir = removeBlanks(c.String())
	c = contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", "export GPCC_INSTANCE_NAME="), []string{"FS", "="})
	i.GPCC.InstanceName = removeBlanks(c.String())
	c = contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", "export GPCC_INSTANCE_PORT="), []string{"FS", "="})
	i.GPCC.InstancePort = removeBlanks(c.String())
	c = contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", "export GPPERFMONHOME="), []string{"FS", "="})
	i.GPCC.GpPerfmonHome = removeBlanks(c.String())
}

// Check if this version of database has GPCC installed.
func (i *Installation) doesThisEnvHasGPCCInstalled() {
	Debugf("Checking if there is already a version of gpcc installed on this database version")
	// If exists then is there a GPCC already installed, if yes then ask for confirmation
	if i.GPCC.InstanceName != "" {
		Warnf("Found a instance of GPCC already installed on this environment, please confirm if the existing GPCC can be uninstalled")
		// Now ask for the confirmation
		confirm := YesOrNoConfirmation()

		// What was the confirmation
		if confirm == "y" { // yes, then uninstall the old GPCC installation
			i.uninstallGPCC()
		} else { // no then exit
			Infof("Cancelling the installation...")
			os.Exit(0)
		}
	}
}

// Install the product that is requested
func (i *Installation) installGPCCProduct() {

	Infof("Installing GPCC version \"%s\" on the GPDB Version \"%s\"", cmdOptions.CCVersion, cmdOptions.Version)
	
}



// Install the product that is requested
func (i *Installation) postGPCCInstall() {

	Infof("Running the post steps for the installation of GPCC version \"%s\" on the GPDB Version \"%s\"", cmdOptions.CCVersion, cmdOptions.Version)

}