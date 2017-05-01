package install

import (
	log "../../pkg/core/logger"
)

import (
	"../../pkg/core/arguments"
	"../../pkg/core/methods"
	"../../pkg/install/library"
	"strings"
	"errors"
)

func UnzipBinary() error {

	// List all the files in the download directory
	all_files, err := library.ListFilesInDir(arguments.DownloadDir)
	if err != nil { return err }

	// Check if any of the file matches to requested version to install
	binary_file, err := library.GetBinaryOfMatchingVersion(all_files, arguments.RequestedInstallVersion)
	if err != nil { return err }

	// If we cannot find a match then error out
	if methods.IsValueEmpty(binary_file) {
		return errors.New("Cannot find any binaries that matches the version: " + arguments.RequestedInstallVersion)
	}

	// Check if the file is a zip file found or Unzip and do work accordingly
	if strings.HasSuffix(binary_file, ".zip") {
		err := library.Unzip(binary_file, arguments.DownloadDir)
		return err
	} else {
		log.Println("Using GPDB binaries found in the download directory: " + binary_file)
		return nil
	}

	return nil
}
