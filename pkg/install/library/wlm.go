package library

import (
	log "../../core/logger"
	"../../core/arguments"
	"../objects"
	"io/ioutil"
	"regexp"
	"strings"
	"../../core/methods"
)

func findWLMBinaries(path string, keyword string, info_text string, warn_text string) (string, error) {

	// Listing all the files in the directory
	output, err := ioutil.ReadDir(path)
	if err != nil { return "", err }
	for _, f := range output {

		// Hunting for file that matches the condition
		re := regexp.MustCompile(``+ keyword +``)
		matches := re.FindStringSubmatch(f.Name())

		// Found the file, return the file name
		if len(matches) != 0 {
			log.Println(info_text + matches[0])
			return matches[0], nil
		}
	}

	// Did we find the file.
	if !methods.IsValueEmpty(warn_text) {
		log.Warn(warn_text)
	} else {
		log.Println(warn_text)
	}

	return "", nil
}

func InstallWLM() error {

	log.Println("Installing workload manager")

	// WLM Binary file name
	log.Println("Checking the directory \"" + objects.GPPERFMONHOME + "\" for WLM binaries.")
	WLMBinaryFile, err := findWLMBinaries(objects.GPPERFMONHOME,
		"gp-wlm.*.bin",
		"Found the WLM Binary file: ",
		"Cannot find the WLM Binaries on the directory \"" + objects.GPPERFMONHOME + "\", skipping the WLM Installation")
	if err != nil { return err }

	// Version and Directory
	WLMVersion := strings.Replace(WLMBinaryFile, ".bin", "", 1)
	WLMDir := arguments.EnvYAML.Install.MasterDataDirectory + "wlm/" + WLMVersion

	// If no binaries are found then return or exit from this function.
	if WLMBinaryFile == "" {
		return nil
	}

	// Check if the binary we are installing is already installed on this server
        HasThisWLMAlreadyInstalled, err := findWLMBinaries( arguments.EnvYAML.Install.MasterDataDirectory + "wlm/",
	        WLMVersion,
	        "Found the WLM Version already installed: ",
	        "Cannot find any previous installation of WLM version \""+ WLMVersion +"\" installed on this server")

	// Yes found a previous installation, so lets skip the installation.
	if !methods.IsValueEmpty(HasThisWLMAlreadyInstalled) {
		return nil
	}

	// Installing the workload manager.
	log.Println("Creating script to install the workload manager")
	var InstallWLMArg []string

	InstallWLMArg = append(InstallWLMArg, "source " + objects.EnvFileName)
	InstallWLMArg = append(InstallWLMArg, "chmod +x $GPPERFMONHOME/" + WLMBinaryFile)
	InstallWLMArg = append(InstallWLMArg, "mkdir -p " + WLMDir)
	InstallWLMArg = append(InstallWLMArg, "$GPPERFMONHOME/" + WLMBinaryFile + " --install=" + WLMDir)

	// Write it to the file and execute.
	file := arguments.TempDir + "install_wlm.sh"
	err = ExecuteBash(file, InstallWLMArg)
	if err != nil { return err }

	// Update the environment file
	err = UpdateWlmEnvFile(WLMDir, WLMVersion)
	if err != nil { return err }

	log.Println("Installation of WLM manager is complete")

	return nil
}


func StopWLMService() error {

	log.Println("Stopping all the workload manager services running on the cluster")

	// Stopping all the WLM arguments.
	var StopWLMArgs []string
	StopWLMArgs = append(StopWLMArgs, objects.WLMInstallDir + "/gp-wlm/bin/svc-mgr.sh --service=all --action=cluster-stop")

	// Write it to the file.
	file := arguments.TempDir + "stop_wlm.sh"
	err := ExecuteBash(file, StopWLMArgs)
	if err != nil { return err }

	return nil
}


func UninstallWLM(t string) error {

	log.Println("Found a old version of WLM: " + objects.WLMVersion)

	// Stop all WLM services
	err := StopWLMService()
	if err != nil { return err }

	// Uninstall Arguments
	var UninstallWLMArgs []string
	UninstallWLMArgs = append(UninstallWLMArgs, "source " + objects.EnvFileName)
	UninstallWLMArgs = append(UninstallWLMArgs, objects.WLMInstallDir + "/gp-wlm/bin/uninstall --symlink " + objects.WLMInstallDir + "/gp-wlm")
	UninstallWLMArgs = append(UninstallWLMArgs, "cp " + objects.EnvFileName + " " + objects.EnvFileName + "." + t)
	UninstallWLMArgs = append(UninstallWLMArgs, "egrep -v \"wlm|WLM\" " + objects.EnvFileName + "." + t + " > " + objects.EnvFileName)
	UninstallWLMArgs = append(UninstallWLMArgs, "rm " + objects.EnvFileName + "." + t)

	// Write it to the file.
	file := arguments.TempDir + "uninstall_wlm.sh"
	err = ExecuteBash(file, UninstallWLMArgs)
	if err != nil { return err }

	return nil
}