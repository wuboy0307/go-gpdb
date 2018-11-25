package main

import (
	"strings"
	"fmt"
)

// Build uninstall script
func (i *Installation) createUninstallScript() error {

	// Uninstall script location
	uninstallFile := Config.INSTALL.UNINTSALLDIR + "uninstall_" + cmdOptions.Version + "_" + i.Timestamp
	Infof("Creating Uninstall file for this installation at: " + uninstallFile)

	// Query
	queryString := `
select $$ssh $$ || hostname || $$ "ps -ef|grep postgres|grep -v grep|grep $$ ||  port || $$ | awk '{print $2}'| xargs -n1 /bin/kill -11 &>/dev/null" $$ from gp_segment_configuration 
union
select $$ssh $$ || hostname || $$ "rm -rf /tmp/.s.PGSQL.$$ || port || $$*"$$ from gp_segment_configuration
union
select $$ssh $$ || c.hostname || $$ "rm -rf $$ || f.fselocation || $$"$$ from pg_filespace_entry f, gp_segment_configuration c where c.dbid = f.fsedbid
`

	// Execute the query
	cmdOut, err := executeOsCommandOutput("psql", "-p", i.GPInitSystem.MasterPort, "-d", "template1", "-Atc", queryString)
	if err != nil {
		Fatalf("Error in running uninstall command on database, err: %v", err)
	}

	// Create the file
	createFile(uninstallFile)
	writeFile(uninstallFile, []string{
		string(cmdOut),
		"rm -rf "+ Config.INSTALL.ENVDIR +"env_" + cmdOptions.Version + "_"+ i.Timestamp,
		"rm -rf " + uninstallFile,

	})
	return nil
}

// Uninstall using gpdeletesystem
func removeEnvGpDeleteSystem(envFile string) error {

	Info("Starting the database if stopped to run the gpdeletesystem on the environment")

	// Start the database if not started
	startDBifNotStarted(envFile)

	Infof("Calling gpdeletesystem to remove the environment: %s", envFile)

	// Write it to the file.
	file := Config.CORE.TEMPDIR + "run_deletesystem.sh"
	createFile(file)
	writeFile(file, []string{
		"source " + envFile,
		"gpdeletesystem -d $MASTER_DATA_DIRECTORY -f << EOF",
		"y",
		"y",
		"EOF",
	})
	_, err := executeOsCommandOutput("/bin/sh", file)
	if err != nil {
		deleteFile(file)
		return err
	}
	deleteFile(file)
	return nil
}

// Uninstall using manual method
func removeEnvManually(version, timestamp string) {
	uninstallScript := Config.INSTALL.UNINTSALLDIR + fmt.Sprintf("uninstall_%s_%s", version, timestamp)
	Infof("Cleaning up the extra files using the uninstall script: %s", uninstallScript)
	exists, err := doesFileOrDirExists(uninstallScript)
	if err != nil {
		Fatalf("error when trying to find the uninstaller file \"%s\", err: %v", uninstallScript, err)
	}
	if exists {
		executeOsCommandOutput("/bin/sh", uninstallScript)
	} else {
		Fatalf("Unable to find the uninstaller file \"%s\"", uninstallScript)
	}
}

// Main Remove method
func remove() {

	Infof("Starting program to uninstall the version: %s", cmdOptions.Version)

	// Check if the envfile for that version exists
	chosenEnvFile := installedEnvFiles(fmt.Sprintf("*%s*", cmdOptions.Version), "choose", true)

	// If we receive none, then display the error to user
	var timestamp, version string
	if IsValueEmpty(chosenEnvFile) {
		Fatalf("Cannot find any environment with the version: %s", cmdOptions.Version)
	} else { // Else store the value
		timestamp = strings.Split(chosenEnvFile, "_")[2]
		version = strings.Split(chosenEnvFile, "_")[1]
	}
	Infof("The choosen enviornment file to remove is: %s ", chosenEnvFile)
	Info("Uninstalling the environment")

	// If there is failure in gpstart, user can use force to force manual uninstallation
	if !cmdOptions.Force {
		err := removeEnvGpDeleteSystem(chosenEnvFile)
		if err != nil {
			Warnf("Failed to uninstall using gpdeletesystem, trying manual method..")
		}
	} else {
		Infof("Forcing uninstall of the environment: %s", chosenEnvFile)
	}

	// Run this to cleanup the file created by go-gpdb
	removeEnvManually(version, timestamp)

	// TODO: Uninstall GPCC

	Infof("Uninstallation of environment \"%s\" was a success", chosenEnvFile)
	Info("exiting ....")
}

