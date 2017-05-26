package library

import (
	log "../../core/logger"
	"../objects"
	"../../core/arguments"
	"../../core/methods"
	"os/exec"
	"errors"
	"strconv"
	"strings"
)


// Installed GPPERFMON software to the database
func InstallGpperfmon() error {

	log.Println("Installing GPPERFMON Database..")

	// Check if the password for gpmon is provided on config.yml, else set it to default "changeme"
	if methods.IsValueEmpty(arguments.EnvYAML.Install.GpMonPass) {
		log.Warn("The value of GPMON_PASS is missing, setting the password of gpperfmon & gpmon defaulted to \"changeme\"")
		arguments.EnvYAML.Install.GpMonPass = "changeme"
	}

	// Installation script
	var install_script []string
	install_script = append(install_script, "source " + objects.EnvFileName)
	install_script = append(install_script, "gpperfmon_install --enable --password "+ arguments.EnvYAML.Install.GpMonPass +" --port $PGPORT")
	install_script = append(install_script, "echo \"host      all     gpmon    0.0.0.0/0    md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf")
	install_script = append(install_script, "echo \"host     all         gpmon         ::1/128       md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf")

	// Execute the script
	temp_file := arguments.TempDir + "install_gpperfmon.sh"
	err := ExecuteBash(temp_file, install_script)
	if err != nil { return err } else {}

	return nil
}

// Check if the GPCC has installed correctly
func WasGPCCInstallationSucess() error {

	// Is gpmmon running
	log.Println("Checking if GPMMON process is running")
	_ , err := exec.Command("pgrep", "gpmmon").Output()
	if err != nil {
		return errors.New("GPMMON process is not running")
	} else {
		log.Println("GPMMON process is running, check if the gpperfmon database has all the required tables.")
	}

	// Can we access gpperfmon database
	log.Println("Checking if GPPERFMON database can be accessed")
	query_string := "select * from system_now"
	_, err = execute_db_query(query_string, "gpperfmon",false, "")
	if err != nil {
		return errors.New("Cannot read the gpperfmon database, installation of gpperfmon database failed..")
	} else {
		log.Println("GPPERFMON Database can be accessed, continuing the script...")
	}

	//

	return nil
}

// Install the Command Center Web UI
func InstallGPCCUI(args []string, cc_home string) error {

	log.Println("Installing command center WEB UI")
	// Check if the version of command center is of 1.x, since that version didn't have WLM.
	var gpcmdrArgs []string
	gpcmdrArgs = append(gpcmdrArgs, "source " + objects.EnvFileName)
	gpcmdrArgs = append(gpcmdrArgs, "source " + cc_home + "/gpcc_path.sh")
	gpcmdrArgs = append(gpcmdrArgs, "echo")
	gpcmdrArgs = append(gpcmdrArgs, "gpcmdr --setup << EOF")
	for _, arg := range args {
		gpcmdrArgs = append(gpcmdrArgs, arg)
	}
	gpcmdrArgs = append(gpcmdrArgs, "echo")

	// Write it to the file and execute
	file := arguments.TempDir + "gpcmdr_setup.sh"
	err := ExecuteBash(file, gpcmdrArgs)
	if err != nil { return err }

	return nil
}

// Install GPCC WebUI
func InstallGPCCWEBUI(cc_name string, ccp int) error {

	log.Println("Running the setup for installing GPCC WEB UI")

	// CC OPtion for different version of cc installer
	if strings.HasPrefix(arguments.RequestedCCInstallVersion, "1") { // CC Version 1.x doesn't have WLM and the option used are
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(objects.ThisDBMasterPort), objects.GPCC_PORT, "n", "n", "n", "n", "EOF"}
		InstallGPCCUI(script_option, objects.BinaryInstallLocation)
	} else if strings.HasPrefix(arguments.RequestedCCInstallVersion, "2.5") { // Option of CC 2.5 & after
		objects.InstallWLM = true
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(objects.ThisDBMasterPort), objects.GPCC_PORT, strconv.Itoa(ccp + 1), "n", "n", "n", "n", "EOF"}
		InstallGPCCUI(script_option, objects.BinaryInstallLocation)
	} else if strings.HasPrefix(arguments.RequestedCCInstallVersion, "2.1") || strings.HasPrefix(arguments.RequestedCCInstallVersion, "2.0") { // Option of CC 2.0 & after
		objects.InstallWLM = true
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(objects.ThisDBMasterPort), "n", objects.GPCC_PORT, "n", "n", "n", "n", "EOF"}
		InstallGPCCUI(script_option, objects.BinaryInstallLocation)
	} else if strings.HasPrefix(arguments.RequestedCCInstallVersion, "2") { // Option for other version of cc 2.x
		objects.InstallWLM = true
		var script_option = []string{cc_name, "n", cc_name, strconv.Itoa(objects.ThisDBMasterPort), "n", objects.GPCC_PORT, strconv.Itoa(ccp + 1), "n", "n", "n", "n", "EOF"}
		InstallGPCCUI(script_option, objects.BinaryInstallLocation)
	} else if strings.HasPrefix(arguments.RequestedCCInstallVersion, "3.0") { // Option for CC version 3.0
		objects.InstallWLM = true
		var script_option = []string{cc_name, cc_name, "n", strconv.Itoa(objects.ThisDBMasterPort), objects.GPCC_PORT, "n", "n", "EOF"}
		InstallGPCCUI(script_option, objects.BinaryInstallLocation)
	} else { // All the newer version option unless changed.
		objects.InstallWLM = true
		var script_option = []string{cc_name, cc_name, "n", strconv.Itoa(objects.ThisDBMasterPort), "n", objects.GPCC_PORT, "n", "n", "EOF"}
		InstallGPCCUI(script_option, objects.BinaryInstallLocation)
	}

	return nil
}

// Store the last used port
func StoreLastUsedGPCCPort() error {

	// Fully qualified filename
	FurtureRefFile := arguments.FutureRefDir + objects.GPCCPortBaseFileName
	log.Println("Storing the last used ports for this installation at: " + FurtureRefFile)

	// Delete the file if already exists
	_ = methods.DeleteFile(FurtureRefFile)
	err := methods.CreateFile(FurtureRefFile)
	if err != nil {return err}

	// Contents to write to the file
	var FutureContents []string
	port, _ := strconv.Atoi(objects.GPCC_PORT)
	port = port + 3
	FutureContents = append(FutureContents, "GPCC_PORT: " + strconv.Itoa(port))

	// Write to the file
	err = methods.WriteFile(FurtureRefFile, FutureContents)
	if err != nil {return err}

	return nil
}

// Uninstall gpcc
func UninstallGPCC(t string, env_file string) error {

	log.Println("Uninstalling the version of command center that is currently installed on this environment.")

	// Generate the script
	var uninstallGPCCArgs []string
	uninstallGPCCArgs = append(uninstallGPCCArgs, "source " + env_file)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "source " + objects.GPPERFMONHOME + "/gpcc_path.sh")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "gpcmdr --stop " + objects.GPCC_INSTANCE_NAME + " &>/dev/null" )
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf " + objects.GPPERFMONHOME + "/instances/" + objects.GPCC_INSTANCE_NAME)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "gpconfig -c gp_enable_gpperfmon -v off &>/dev/null")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "cp $MASTER_DATA_DIRECTORY/pg_hba.conf $MASTER_DATA_DIRECTORY/pg_hba.conf." + t)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "grep -v gpmon $MASTER_DATA_DIRECTORY/pg_hba.conf."+ t +" > $MASTER_DATA_DIRECTORY/pg_hba.conf")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf $MASTER_DATA_DIRECTORY/pg_hba.conf."+ t )
	uninstallGPCCArgs = append(uninstallGPCCArgs, "psql -d template1 -p $PGPORT -Atc \"drop database gpperfmon\" &>/dev/null")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "psql -d template1 -p $PGPORT -Atc \"drop role gpmon\" &>/dev/null")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf $MASTER_DATA_DIRECTORY/gpperfmon/*")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "cp "+ env_file +" "+ env_file + "." + t )
	uninstallGPCCArgs = append(uninstallGPCCArgs, "egrep -v \"GPPERFMONHOME|GPCC_INSTANCE_NAME|GPCCPORT\" " + env_file + "." + t + " > " + env_file)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf "+ env_file + "." + t)

	// Write it to the file.
	file := arguments.TempDir + "uninstall_gpcc.sh"
	err := ExecuteBash(file, uninstallGPCCArgs)
	if err != nil { return err }

	return nil
}
