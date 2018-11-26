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

	// The installation b/w GPCC 3 & GPCC 4 is completely different so
	// we need to create two working libraries to make it work.
	i.WorkingHostFileLocation = os.Getenv("HOME") + "/hostfile"
	exists, _  := doesFileOrDirExists(i.WorkingHostFileLocation)
	if !exists {
		Fatalf("No host file found on the home location of the user")
	}
	if isThis4x() {
		i.installGPCC4xAndAbove()
	} else {
		i.installGPCCBelow4x()
	}

}

// Is it a 4.x to <4.x version
func isThis4x() bool {
	v := extractVersion(cmdOptions.CCVersion)
	if v > 4.0 { // newer GPCC code
		return true
	} else { // legacy GPCC code
		return false
	}
}

// Install GPCC 4.x
func (i *Installation) installGPCC4xAndAbove() {
	// unzip the binaries
	i.GPCC.GPCCBinaryLoc = unzip(fmt.Sprintf("*%s*", cmdOptions.CCVersion))

	// Install the Gpperfmon database
	i.installGpperfmon()

	// Install the gpcc binaries.
	i.installGPCCbinaries4x()

	// TODO: Start GPCC
}

// Install GPCC less than 4.x
func (i *Installation) installGPCCBelow4x() {
	i.GPCC.GPCCBinaryLoc = unzip(fmt.Sprintf("*%s*", cmdOptions.CCVersion))
}

// Install the product that is requested
func (i *Installation) postGPCCInstall() {

	Infof("Running the post steps for the installation of GPCC version \"%s\" on the GPDB Version \"%s\"", cmdOptions.CCVersion, cmdOptions.Version)

	// Store the last used port for future use
	i.savePort()

	// Update the envfile
	i.updateEnvFile()

	// Installation complete, print on the screen the env file to source and cleanup temp files
	displayEnvFileToSource(i.EnvFile)

	// get the ip of the host and display on the screen
	ip, _ := GetLocalIP()
	Infof("GPCC Web URL: http://%s:%s, username / password: gpmon / %s", ip, i.GPCC.InstancePort, Config.INSTALL.GPMONPASS)
}

// Installed GPPERFMON software to the database
func (i *Installation) installGpperfmon() error {

	Infof("Installing GPPERFMON Database on the database with environment file: %s", i.EnvFile)

	// Installation script
	GpPerfMonFile := Config.CORE.TEMPDIR + "install_gpperfmon.sh"
	generateBashFileAndExecuteTheBashFile(GpPerfMonFile, "/bin/sh", []string{
		"source " + i.EnvFile,
		"gpperfmon_install --enable --password " + Config.INSTALL.GPMONPASS + " --port $PGPORT",
		"echo \"local    gpperfmon   gpmon         md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf",
		"echo \"host      all     gpmon    0.0.0.0/0    md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf",
		"echo \"host     all         gpmon         ::1/128       md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf",
	})

	// Verify if the gpperfmon process are running
	i.verifyGpperfmon()
	return nil
}

// Verify if database gpperfmon is installed correctly
func (i *Installation) verifyGpperfmon() {

	// Restart
	Infof("Restarting the database so that the command center agents can be started")
	stopDB(i.EnvFile)
	startDB(i.EnvFile)

	// Verify GPmon process
	Infof("Checking if GPMMON process is running")
	_, err := executeOsCommandOutput("pgrep", "gpmmon")
	if err != nil {
		Fatalf("GPPMON process is not running, the installation of gpperfmon was a failure, err: %$v", err)
	} else {
		Infof("GPMMON process is running, check if the gpperfmon database has all the required tables.")
	}

	// Verify gpperfmon database
	Infof("Checking if GPPERFMON database can be accessed")
	queryString := "select * from system_now"
	gpperfmonHealthFile := Config.CORE.TEMPDIR + "gpperfmon_health.sh"
	deleteFile(gpperfmonHealthFile)
	createFile(gpperfmonHealthFile)
	writeFile(gpperfmonHealthFile, []string{
		fmt.Sprintf("source %s", i.EnvFile),
		fmt.Sprintf("psql -d gpperfmon -U gpmon -c '%s'", queryString),
	})
	_, err = executeOsCommandOutput("/bin/sh", gpperfmonHealthFile)
	if err != nil {
		Fatalf("GPPERFMON database is not accessible, err: %v", err)
	} else {
		Infof("GPPERFMON Database can be accessed, continuing the script...")
	}
}

// Create configuration file
func (i *Installation) createGPCCCOnfigurationFile() string {
	i.GPCC.InstancePort = i.validatePort( "GPCC_PORT", defaultGpccPort)
	i.GPCC.WebSocketPort = i.validatePort( "WEBSOCKET_PORT", defaultWebSocket)
	i.GPCC.InstanceName = fmt.Sprintf("ccversion_%s", cmdOptions.CCVersion)
	gpccInstallConfig := Config.CORE.TEMPDIR + fmt.Sprintf("gpcc_config_4x_%s.sh", i.Timestamp)
	deleteFile(gpccInstallConfig)
	createFile(gpccInstallConfig)
	writeFile(gpccInstallConfig, []string{
		"path = /usr/local",
		"display_name = " + i.GPCC.InstanceName,
		"master_port = " + i.GPInitSystem.MasterPort,
		"web_port = " + i.GPCC.InstancePort,
		"rpc_port = " + i.GPCC.WebSocketPort,
		"enable_ssl = false",
		"enable_kerberos = false",
	})
	return gpccInstallConfig
}

// Installing gpcc instance
func (i *Installation) installGPCCbinaries4x() {
	Infof("Installing gpcc binaries for the cc version: %s", cmdOptions.CCVersion)
	configFile := i.createGPCCCOnfigurationFile()
	executeOsCommand(i.GPCC.GPCCBinaryLoc, "-c", configFile, "y", "EOF", "&> /dev/null")
}