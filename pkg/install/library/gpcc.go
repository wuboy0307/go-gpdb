package library

import (
	log "../../core/logger"
	"../objects"
	"../../core/arguments"
	"../../core/methods"
	"os/exec"
	"errors"
	"strconv"
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

// Started GPCC WebUI
func StartGPCC (cc_name string, cc_home string) error {

	log.Println("Starting command center WEB UI")

	// Start the command center web UI
	var gpcmdrStartarg []string
	gpcmdrStartarg = append(gpcmdrStartarg, "source " + objects.EnvFileName)
	gpcmdrStartarg = append(gpcmdrStartarg, "source " + cc_home + "/gpcc_path.sh")
	gpcmdrStartarg = append(gpcmdrStartarg, "echo")
	gpcmdrStartarg = append(gpcmdrStartarg, "gpcmdr --start " + cc_name)
	gpcmdrStartarg = append(gpcmdrStartarg, "echo")

	// Write it to the file.
	file := arguments.TempDir + "gpcmdr_start.sh"
	err := ExecuteBash(file, gpcmdrStartarg)
	if err != nil { return err }

	return nil
}


// Update Environment file
func UpdateEnvFile(cc_name string) error {

	log.Println("Updating the environment file \"" + objects.EnvFileName + "\" with the GPCC environment")

	// Environment file contents
	var EnvFileContents []string
	EnvFileContents = append(EnvFileContents, "export GPPERFMONHOME=" + objects.BinaryInstallLocation)
	EnvFileContents = append(EnvFileContents, "export PATH=$GPPERFMONHOME/bin:$PATH")
	EnvFileContents = append(EnvFileContents, "export LD_LIBRARY_PATH=$GPPERFMONHOME/lib:$LD_LIBRARY_PATH")
	EnvFileContents = append(EnvFileContents, "export GPCC_INSTANCE_NAME=" + cc_name)
	EnvFileContents = append(EnvFileContents, "export GPCCPORT=" + objects.GPCC_PORT)

	// Append to file
	err := methods.AppendFile(objects.EnvFileName, EnvFileContents)
	if err != nil {return nil}

	return nil
}

// Store the last used port
func StoreLastUsedGPCCPort() error {

	// Fully qualified filename
	FurtureRefFile := arguments.FutureRefDir + objects.GPCCPortBaseFileName
	log.Println("Storing the last used ports for this installation at: " + FurtureRefFile)

	// Delete the file if already exists
	err := methods.DeleteFile(FurtureRefFile)
	if err != nil {return nil}
	err = methods.CreateFile(FurtureRefFile)
	if err != nil {return nil}

	// Contents to write to the file
	var FutureContents []string
	port, _ := strconv.Atoi(objects.GPCC_PORT)
	port = port + 1
	FutureContents = append(FutureContents, "GPCC_PORT: " + strconv.Itoa(port))

	// Write to the file
	err = methods.WriteFile(FurtureRefFile, FutureContents)
	if err != nil {return nil}

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
	uninstallGPCCArgs = append(uninstallGPCCArgs, "psql -d template1 -p $PGPORT -Atc \"drop database if exists gpperfmon\" &>/dev/null")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "psql -d template1 -p $PGPORT -Atc \"drop role gpmon\" &>/dev/null")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf $MASTER_DATA_DIRECTORY/gpperfmon/*")
	uninstallGPCCArgs = append(uninstallGPCCArgs, "cp "+ env_file +" "+ env_file + "." + t )
	uninstallGPCCArgs = append(uninstallGPCCArgs, "egrep -v \"GPPERFMONHOME|GPCC_INSTANCE_NAME|GPCCPORT\" " + env_file + "." + t + " > " + env_file)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "rm -rf "+ env_file + "." + t)
	uninstallGPCCArgs = append(uninstallGPCCArgs, "gpstop -afr")

	// Write it to the file.
	file := arguments.TempDir + "uninstall_gpcc.sh"
	err := ExecuteBash(file, uninstallGPCCArgs)
	if err != nil { return err }

	return nil
}