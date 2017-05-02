package install


import (
	"../../pkg/core/arguments"
	"../../pkg/core/methods"
	"../../pkg/install/library"
	"../../pkg/install/objects"
)


func InstallSingleNodeGPDB() error {

	// If the install is called from download command the set default values
	if !methods.IsValueEmpty(arguments.RequestedDownloadVersion) {
		arguments.RequestedInstallVersion = arguments.RequestedDownloadVersion
	}

	// Unzip the binaries, if its file is zipped
	binary_file, err := UnzipBinary()
	if err != nil { return err }

	// Check if there is already a previous version of the same version


	// execute the binaries.
	binary_installation_loc := "/usr/local/greenplum-db-" + arguments.RequestedInstallVersion
	var script_option = []string{"yes", binary_installation_loc, "yes", "yes"}
	err = library.ExecuteBinaries(binary_file, objects.InstallGPDBBashFileName, script_option)
	if err != nil { return err }

	// Generate gpinitsystem config file

	// Execute gpinitsystem

	// Create Environment file for this installation

	// Uninstall script for this installation

	// Store that last used port


	return nil

}
