package library

import (
	log "../../core/logger"
	"../../core/arguments"
	"../../core/methods"
	"../objects"
	"strconv"
	"io/ioutil"
	"strings"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"errors"
)

// Create environment file of this installation
func CreateEnvFile(t string) error {

	// Environment file fully qualified path
	objects.EnvFileName = arguments.EnvFileDir + "env_" + arguments.RequestedInstallVersion + "_" + t
	log.Println("Creating environment file for this installation at: " + objects.EnvFileName)

	// Create the file
	err := methods.CreateFile(objects.EnvFileName)
	if err != nil { return err }

	// Build arguments to write
	var EnvFileContents []string
	EnvFileContents = append(EnvFileContents, "export GPHOME=" + objects.BinaryInstallLocation)
	EnvFileContents = append(EnvFileContents, "export PYTHONPATH=$GPHOME/lib/python")
	EnvFileContents = append(EnvFileContents, "export PYTHONHOME=$GPHOME/ext/python")
	EnvFileContents = append(EnvFileContents, "export PATH=$GPHOME/bin:$PYTHONHOME/bin:$PATH")
	EnvFileContents = append(EnvFileContents, "export LD_LIBRARY_PATH=$GPHOME/lib:$PYTHONHOME/lib:$LD_LIBRARY_PATH")
	EnvFileContents = append(EnvFileContents, "export OPENSSL_CONF=$GPHOME/etc/openssl.cnf")
	EnvFileContents = append(EnvFileContents, "export MASTER_DATA_DIRECTORY=" + objects.GpInitSystemConfig.MasterDir + "/" + objects.GpInitSystemConfig.ArrayName + "-1")
	EnvFileContents = append(EnvFileContents, "export PGPORT=" + strconv.Itoa(objects.GpInitSystemConfig.MasterPort))
	EnvFileContents = append(EnvFileContents, "export PGDATABASE=" + objects.GpInitSystemConfig.DBName)

	// Write to EnvFile
	err = methods.WriteFile(objects.EnvFileName, EnvFileContents)
	if err != nil { return err }

	return nil
}

// Obtain all the files in the environment directory
func ListEnvFile(MatchingFilesInDir []string) ([]string, error) {

	// Show all the environment files
	log.Println("Found below matching environment file for the version: " + arguments.RequestedInstallVersion)

	// Temp files
	temp_env_file := arguments.TempDir + "temp_env.sh"
	temp_env_out_file := arguments.TempDir + "temp_env.out"

	// Create those files
	_ = methods.DeleteFile(temp_env_file)
	err := methods.CreateFile(temp_env_file)
	if err != nil { return []string{}, err }

	// Get the IP address from the server
	ip, err := GetLocalIP()
	if err != nil { return []string{}, err }

	// Bash script
	var cmd []string

	bashCmd := "incrementor=1" +
		";echo -e \"\nID\tMaster Port\tStatus \t\tEnvironment File\t\t     GPCC Instance Name (GPCC URL)\"   > " + temp_env_out_file +
		";echo \"--------------------------------------------------------------------------------------------------------------------------------------------------\"    >> " + temp_env_out_file +
		";ls -1 " + arguments.EnvFileDir + " | grep env_"+ arguments.RequestedInstallVersion +" | while read line" +
		";do    " +
		"       unset GPCC_INSTANCE_NAME; unset GPCCPORT; unset WLM_VERSION ; unset WLM_PATH" +
		"       ;WLM_PATH=`grep "+ arguments.EnvYAML.Install.MasterDataDirectory + "wlm/ " + arguments.EnvFileDir +"${line} | grep source | awk '{print $2}'`" +
		"       ;if [ ! -f $WLM_PATH ]; then sed -i -e '/(^source.*wlm.*|.*WLM.*=.*)/d' "+ arguments.EnvFileDir + "$line ;fi" +
		"       ;source "+arguments.EnvFileDir+"$line" +
		"       ;psql -d template1 -p $PGPORT -Atc \"select 1\" &>/dev/null" +
		"       ;retcode=$?" +
		"       ;if [ \"$retcode\" == \"0\" ]; then" +
		"               echo -e \"$incrementor\t$PGPORT\t\tRUNNING\t\t$line\t   $GPCC_INSTANCE_NAME (http://"+ ip +":$GPCCPORT)\" >> " + temp_env_out_file +
		"       ;else" +
		"               echo -e \"$incrementor\t$PGPORT\t\tSTOPPED\t\t$line\t   $GPCC_INSTANCE_NAME (http://"+ ip +":$GPCCPORT)\"  >> " + temp_env_out_file +
		"       ;fi" +
		"       ;incrementor=$((incrementor+1))" +
		";done"

	cmd = append(cmd, bashCmd)

	// Copy it to the file
	_ = methods.WriteFile(temp_env_file, cmd)

	// Execute the script
	_, err = exec.Command("/bin/sh", temp_env_file).Output()
	if err != nil { return []string{}, err }

	// Display the output
	out, _ := ioutil.ReadFile(temp_env_out_file)
	outReplace := strings.Replace(string(out), "(http://"+ ip +":)", "", -1)
	fmt.Println(outReplace)

	// Cleanup the temp files
	_ = methods.DeleteFile(temp_env_file)
	_ = methods.DeleteFile(temp_env_out_file)

	// Create a list of the options
	var envStore []string
	for _, e := range MatchingFilesInDir {
		envStore = append(envStore, e)
	}

	return envStore, nil

}

// Check if there is any previous installation of the same version
func PrevEnvFile(product string) (string, error) {

	log.Println("Checking if there is previous installation for the version: " + arguments.RequestedInstallVersion)

	// Read all the file in the env directory
	var MatchingFilesInDir []string
	allfiles, err := ioutil.ReadDir(arguments.EnvFileDir)
	if err != nil { return "", err }

	// Match the version on this directory
	for _, file := range allfiles {

		if strings.Contains(file.Name(), arguments.RequestedInstallVersion) {
			MatchingFilesInDir = append(MatchingFilesInDir, file.Name())
		}

	}

	// Found matching environment file of this installation, now ask for confirmation
	if len(MatchingFilesInDir) > 0 && product == "confirm" {

		_, err := ListEnvFile(MatchingFilesInDir)
		if err != nil { return "", err }

		// Now ask for the confirmation
		confirm := methods.YesOrNoConfirmation()

		// What was the confirmation
		if confirm == "y" {  // yes
			log.Println("Continuing with the installtion of version: " + arguments.RequestedInstallVersion)
			return  MatchingFilesInDir[0], nil
		} else { // no
			log.Println("Cancelling the installation...")
			os.Exit(0)
		}

	} else if len(MatchingFilesInDir) > 1 && product == "choose" { // else choose

		envStore, err := ListEnvFile(MatchingFilesInDir)
		if err != nil { return "", err }

		// What is users choice
		choice := methods.Prompt_choice(len(envStore))

		// return the enviornment file to the main function
		choosenEnv := envStore[choice-1]
		return choosenEnv, nil

	} else if len(MatchingFilesInDir) == 1 && product == "choose" { // if there is only one , then there is no choose just provide the only one

		return MatchingFilesInDir[0], nil

	} else if product == "list" {
		_, err := ListEnvFile(MatchingFilesInDir)
		if err != nil { return "", err }
	}

	// return the environment file.
	return "", err
}

// Set Environment of the shell
func SetVersionEnv(filename string) error {

	log.Println("Attempting to open a terminal, after setting the environment of this installation.")

	// User Home
	usersHomeDir := os.Getenv("HOME")

	// Create a temp file to execute
	executeFile := arguments.TempDir + "openterminal.sh"
	_ = methods.DeleteFile(executeFile)
	_ = methods.CreateFile(executeFile)

	// The command
	var cmd []string
	cmdString := "gnome-terminal --working-directory=\"" + usersHomeDir + "\" --tab -e 'bash -c \"echo \\\"Sourcing Envionment file: "+ filename + "\\\"; source "+ filename +"; exec bash\"'"
	cmd = append(cmd, cmdString)

	// Write to the file
	_ = methods.WriteFile(executeFile, cmd)
	_, err := exec.Command("/bin/sh", executeFile).Output()
	if err != nil {
		log.Println("Couldn't open a new terminal, please copy / paste the below environment location to set the environment")
		log.Println("\n source " + filename + "\n")
	}

	// Cleanup the file file.
	_ = methods.DeleteFile(executeFile)

	return nil
}

// grepping for keyword from file
func ExtractKeywordFromFile(env_file string, kwrd string) ([]string, error) {

	content, err := ioutil.ReadFile(env_file)
	if err != nil { return []string{""}, nil }
	re := regexp.MustCompile(``+ kwrd + ``)
	matches := re.FindStringSubmatch(string(content))

	return matches, nil
}


// Extract PORT and GPHOME location
func ExtractEnvVariables(env_file string) error {

	// Read the file and extract the PGPORT
	matches, err := ExtractKeywordFromFile(env_file, ".*PGPORT=.*")
	if err != nil { return err }

	// Check if we find the PGPORT
	if len(matches) == 0 {
		return errors.New("Cannot find PGPORT value in the environment file: " + env_file)
	} else {
		port := strings.Split(matches[0], "=")[1]
		objects.ThisDBMasterPort, err = strconv.Atoi(port)
		if err != nil { return err }
	}

	// Read the file and extract GPHOME
	matches, err = ExtractKeywordFromFile(env_file, ".*GPHOME=.*")
	if err != nil { return err }

	// Check if we find the GPHOME
	if len(matches) == 0 {
		return errors.New("Cannot find GPHOME value in the environment file: " + env_file)
	} else {
		gphome := strings.Split(matches[0], "=")[1]
		objects.BinaryInstallLocation = gphome
	}

	// extract GPPERFMON instance name
	matches, err = ExtractKeywordFromFile(env_file, ".*GPCC_INSTANCE_NAME=.*")
	if err != nil { return err }

	// Check if we find the CC_instance
	if len(matches) != 0 {
		ccinstance := strings.Split(matches[0], "=")[1]
		objects.GPCC_INSTANCE_NAME = ccinstance
	}

	// extract GPCC_PORT
	matches, err = ExtractKeywordFromFile(env_file, ".*GPCCPORT=.*")
	if err != nil { return err }

	// Check if we find the GPCCPORT
	if len(matches) != 0 {
		ccport := strings.Split(matches[0], "=")[1]
		objects.ThisEnvGPCCPort, _ = strconv.Atoi(ccport)
	}

	// extract GPPERFMONHOME
	matches, err = ExtractKeywordFromFile(env_file, ".*GPPERFMONHOME=.*")
	if err != nil { return err }

	// Check if we find the GPPERFMONHOME
	if len(matches) != 0 {
		gpperfhome := strings.Split(matches[0], "=")[1]
		objects.GPPERFMONHOME = gpperfhome
	}

	// Extract WLM Version
	matches, err = ExtractKeywordFromFile(env_file, ".*WLM_VERSION=.*")
	if err != nil { return err }

	// Check if we find the WLM Version
	if len(matches) != 0 {
		wlmversion := strings.Split(matches[0], "=")[1]
		objects.WLMVersion = wlmversion
	}

	// Extract WLM Install Directory
	matches, err = ExtractKeywordFromFile(env_file, ".*WLM_INSTALL_DIR=.*")
	if err != nil { return err }

	// Check if we find the WLM Install directory
	if len(matches) != 0 {
		wlminstalldir := strings.Split(matches[0], "=")[1]
		objects.WLMInstallDir = wlminstalldir
	}

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
	if err != nil {return err}

	return nil
}

// Update WLM to Environment file
func UpdateWlmEnvFile(WLMDir string, WLMVersion string) error {

	log.Println("Updating the environment file \"" + objects.EnvFileName + "\", with workload manager path")

	// extract the version of the wlm
	//wlmv, err := GetWLMVersion(WLMDir)

	// Now update the environment file with the WLM Information.
	var WLMEnvArgs []string
	WLMEnvArgs = append(WLMEnvArgs, "source " + WLMDir + "/gp-wlm/gp-wlm_path.sh")
	WLMEnvArgs = append(WLMEnvArgs, "export WLM_INSTALL_DIR=" + WLMDir)
	WLMEnvArgs = append(WLMEnvArgs, "export WLM_VERSION=" + WLMVersion)

	// Append to file
	err := methods.AppendFile(objects.EnvFileName, WLMEnvArgs)
	if err != nil {return err}

	return nil
}