package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Check if the directory provided is readable and writeable
func dirValidator() error {

	Debugf("Checking for the existences of master and segment directory")

	// Check if the master & segment location is readable and writable
	masterDirExists, err := doesFileOrDirExists(Config.INSTALL.MASTERDATADIRECTORY)
	if err != nil {
		Fatalf("Cannot get the information of directory %s, err: %v", Config.INSTALL.MASTERDATADIRECTORY, err)
	}

	segmentDirExists, err := doesFileOrDirExists(Config.INSTALL.SEGMENTDATADIRECTORY)
	if err != nil {
		Fatalf("Cannot get the information of directory %s, err: %v", Config.INSTALL.MASTERDATADIRECTORY, err)
	}

	// if the file doesn't exists then let try creating it ...
	if !masterDirExists || !segmentDirExists {
		CreateDir(Config.INSTALL.MASTERDATADIRECTORY)
		CreateDir(Config.INSTALL.SEGMENTDATADIRECTORY)
	}

	return nil
}


// Check if the binaries exits and unzip the binaries.
func (i *Installation) getBinaryFile() string {
	return unzip(fmt.Sprintf("*%s*", cmdOptions.Version))
}

// Installing the greenplum binaries
func (i *Installation) installProduct() {
	// Installing the binaries.
	binFile := i.getBinaryFile()
	Infof("Using the bin file to install the GPDB Product: %s", binFile)
	i.BinaryInstallationLocation = "/usr/local/greenplum-db-" + cmdOptions.Version
	var scriptOption = []string{"yes", i.BinaryInstallationLocation, "yes", "yes"}
	err := executeBinaries(binFile, "install_software.sh", scriptOption)
	if err != nil {
		Fatalf("Failed in installing the binaries, err: %v", err)
	}
}

// Check if the provided hostnames are valid
func (i *Installation) CheckHostnameIsValid() error {

	// Get the hostname of the master
	i.GPInitSystem.MasterHostname = os.Getenv("HOSTNAME")
	if i.GPInitSystem.MasterHostname == "" {
		Fatalf("The environment variable 'HOSTNAME' for master host is not set")
	}

	_, err := exec.Command("ssh", i.GPInitSystem.MasterHostname, "-o", "ConnectTimeout=5", "echo 1").Output()
	if err != nil {
		return err
	}

	// Enable passwordless login
	err = ExecuteGpsshExkey()
	if err != nil {
		return err
	}

	return nil
}

// Run keyless access to the server
func ExecuteGpsshExkey() error {

	// Source GPDB PATH
	err := SourceGPDBPath()
	if err != nil {
		return err
	}

	// Checking if the username and password parameters are passed
	// TODO: Clean this up
	if core.IsValueEmpty(core.EnvYAML.Install.MasterUser) {
		return errors.New("MASTER_USER parameter missing in the config file, please set it and try again")
	}
	if core.IsValueEmpty(core.EnvYAML.Install.MasterPass) {
		return errors.New("MASTER_PASS parameter missing in the config file, please set it and try again")
	}

	// Execute gpssh script to enable keyless access
	Infof("Running gpssh-exkeys to enable keyless access on this server")
	cmd := exec.Command(os.Getenv("GPHOME")+"/bin/gpssh-exkeys", "-h", GpInitSystemConfig.MasterHostName)
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

// Source greenplum path
func (i *Installation) sourceGPDBPath() error {

	// Setting up greenplum path
	err := os.Setenv("GPHOME", i.BinaryInstallationLocation)
	if err != nil {
		return err
	}
	err = os.Setenv("PYTHONPATH", os.Getenv("GPHOME")+"/lib/python")
	if err != nil {
		return err
	}
	err = os.Setenv("PYTHONHOME", os.Getenv("GPHOME")+"/ext/python")
	if err != nil {
		return err
	}
	err = os.Setenv("PATH", os.Getenv("GPHOME")+"/bin:"+os.Getenv("PYTHONHOME")+"/bin:"+os.Getenv("PATH"))
	if err != nil {
		return err
	}
	err = os.Setenv("LD_LIBRARY_PATH", os.Getenv("GPHOME")+"/lib:"+os.Getenv("PYTHONHOME")+"/lib:"+os.Getenv("LD_LIBRARY_PATH"))
	if err != nil {
		return err
	}
	err = os.Setenv("OPENSSL_CONF", os.Getenv("GPHOME")+"/etc/openssl.cnf")
	if err != nil {
		return err
	}

	return nil
}
