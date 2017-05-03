package install

import (
	"../../pkg/core/arguments"
	"../../pkg/core/methods"
	"../../pkg/install/library"
	"../../pkg/install/objects"
	"time"
)

func InstallSingleNodeGPDB() error {

	// If the install is called from download command the set default values
	if !methods.IsValueEmpty(arguments.RequestedDownloadVersion) {
		arguments.RequestedInstallVersion = arguments.RequestedDownloadVersion
	}

	// Validate the master & segment exists and is readable
	err := library.DirValidator(arguments.EnvYAML.Install.MasterDataDirectory, arguments.EnvYAML.Install.SegmentDataDirectory)
	if err != nil { return err }

	// Unzip the binaries, if its file is zipped
	binary_file, err := UnzipBinary()
	if err != nil { return err }

	// Check if there is already a previous version of the same version


	// execute the binaries.
	binary_installation_loc := "/usr/local/greenplum-db-" + arguments.RequestedInstallVersion
	var script_option = []string{"yes", binary_installation_loc, "yes", "yes"}
	err = library.ExecuteBinaries(binary_file, objects.InstallGPDBBashFileName, script_option)
	if err != nil { return err }

	// Source greenplum home
	err = library.SourceGPDBPath(binary_installation_loc)  // Sourcing doesn't work since it looses the variable after exit from shell
	if err != nil { return err }

	// Check ssh to host is working and enable password less login
	err = library.CheckHostnameIsValid(binary_installation_loc)
	if err != nil { return err }

	// Generate gpinitsystem config file
	t := time.Now().Format("20060102150405")
	err = library.BuildGpInitSystemConfig(t)
	if err != nil { return err }

	// Stop any database if any


	// Execute gpinitsystem ( For some reason the gpinitsystem produces exit code
	// 1 even then command ran successfully , so we ignore the exit code 1 here
	err = library.ExecuteGpInitSystem()
	if err != nil && err.Error() != "exit status 1" { return err }

	// Check if the database is healthy
	err = library.IsDBHealthy()
	if err != nil { return err }

	// Create Environment file for this installation
	err = library.CreateEnvFile(t, binary_installation_loc)
	if err != nil { return err }

	// Source the environment that is build
	err = library.SourceEnvFile(t)   // Need to check
	if err != nil { return err }

	// Uninstall script for this installation
	err = library.CreateUnistallScript(t)   // Not storing on the file
	if err != nil { return err }

	// Check if the port is not greater than 63000

	// Store that last used port
	err = library.StoreLastUsedPort()  // Not storing the data
	if err != nil { return err }


	return nil

}
