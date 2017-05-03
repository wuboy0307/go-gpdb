package library

import (
	"errors"
	"os"
	"os/exec"
)

import (
	"../../install/objects"
	"../../core/arguments"
	"../../core/methods"
	log "../../core/logger"
)

// Check if the directory provided is readable and writeable
func DirValidator(master string, segment string) error {

	// Check if the master & segment location is provided.
	if methods.IsValueEmpty(master) {
		return errors.New("MASTER_DATA_DIRECTORY parameter missing in the config file, please set it and try again")
	} else {
		objects.MasterDIR = master
	}

	if methods.IsValueEmpty(segment) {
		return errors.New("SEGMENT_DATA_DIRECTORY parameter missing in the config file, please set it and try again")
	} else {
		objects.SegmentDIR = segment
	}

	// Check if the master & segment location is readable and writable
	master_dir, err := methods.DoesFileOrDirExists(objects.MasterDIR)
	if err != nil { return err }

	segment_dir, err := methods.DoesFileOrDirExists(objects.SegmentDIR)
	if err != nil { return err }

	// if the file doesn't exists then let try creating it ...
	if !master_dir || !segment_dir {
		err := os.MkdirAll(objects.MasterDIR, 0775)
		if err != nil { return err }
		err = os.MkdirAll(objects.SegmentDIR, 0775)
		if err != nil { return err }
	}

	return nil
}

// Check if the provided hostnames are valid
func CheckHostnameIsValid(binary_loc string) error {

	// Check Master host parameter is set
	if methods.IsValueEmpty(arguments.EnvYAML.Install.MasterHost) {
		return errors.New("MASTER_HOST parameter missing in the config file, please set it and try again")
	} else {
		objects.GpInitSystemConfig.MasterHostName = arguments.EnvYAML.Install.MasterHost
	}

	// Check if the provided hostname can be ssh'ed
	hostname := objects.GpInitSystemConfig.MasterHostName
	log.Println("Checking connectivity to host \""+ hostname + "\" can be established")

	_, err := exec.Command("ssh", hostname, "-o" , "ConnectTimeout=5", "echo 1").Output()
	if err != nil { return err }

	// Enable passwordless login
	err = ExecuteGpsshExkey(binary_loc)
	if err != nil { return err }

	return nil
}


// Run keyless access to the server
func ExecuteGpsshExkey(binary_loc string) error {

	// Checking if the username and password parameters are passed
	if methods.IsValueEmpty(arguments.EnvYAML.Install.MasterUser) {
		return errors.New("MASTER_USER parameter missing in the config file, please set it and try again")
	}
	if methods.IsValueEmpty(arguments.EnvYAML.Install.MasterPass) {
		return errors.New("MASTER_PASS parameter missing in the config file, please set it and try again")
	}

	// Execute gpssh script to enable keyless access
	log.Println("Running gpssh-exkeys to enable keyless access on this server")
	cmd := exec.Command("gpssh-exkeys", "-h", objects.GpInitSystemConfig.MasterHostName)
	err := cmd.Run()
	if err != nil { return err }

	return nil
}

// Source greenplum path
func SourceGPDBPath(gphome_loc string) error {

	gphome_loc = gphome_loc + "/greenplum_path.sh"
	log.Println("Sourcing the greenplum path: " + gphome_loc)

	// Source the file
	cmdOut := exec.Command("cat", gphome_loc)
	err :=  cmdOut.Run()
	cmdOut.Wait()
	if err != nil { return err }

	return  nil
}