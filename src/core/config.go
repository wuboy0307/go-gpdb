package core

import (
	"io/ioutil"
	"os"
)

import (
	"../gopkg.in/yaml.v2"
	"../../pkg/core/arguments"
	"../../pkg/core/methods"
	log "../../pkg/core/logger"
)

func CreateDir() error {

	log.Println("Checking if the directories needed for the program exists")
	if methods.IsValueEmpty(arguments.EnvYAML.Core.BaseDir) {
		log.Warn("BASE_DIR parameter missing in the config file, setting to default")
		arguments.EnvYAML.Core.BaseDir = "/home/gpadmin/"
	}

	// App name
	if methods.IsValueEmpty(arguments.EnvYAML.Core.AppName) {
		log.Warn("APPLICATION_NAME parameter missing in the config file, setting to default")
		arguments.EnvYAML.Core.AppName = "gpdbinstall"
	}

	// Temp Directory
	if methods.IsValueEmpty(arguments.EnvYAML.Core.TempDir) {
		log.Warn("TEMP_DIR parameter missing in the config file, setting to default")
		arguments.EnvYAML.Core.TempDir = "/temp/"
	}

	// Download Directory
	if methods.IsValueEmpty(arguments.EnvYAML.Download.DownloadDir) {
		log.Warn("DOWNLOAD_DIR parameter missing in the config file, setting to default")
		arguments.EnvYAML.Download.DownloadDir = "/download/"
	}

	// Env Directory
	if methods.IsValueEmpty(arguments.EnvYAML.Install.EnvDir) {
		log.Warn("ENV_DIR parameter missing in the config file, setting to default")
		arguments.EnvYAML.Install.EnvDir = "/env/"
	}

	// Uninstall Directory
	if methods.IsValueEmpty(arguments.EnvYAML.Install.UnistallDir) {
		log.Warn("UNINTSALL_DIR parameter missing in the config file, setting to default")
		arguments.EnvYAML.Install.UnistallDir = "/uninstall/"
	}

	// Future Reference Directory
	if methods.IsValueEmpty(arguments.EnvYAML.Install.FutureRef) {
		log.Warn("UNINTSALL_DIR parameter missing in the config file, setting to default")
		arguments.EnvYAML.Install.FutureRef = "/future_reference/"
	}

	// Check if the directory exists, else create one.
	base_dir := arguments.EnvYAML.Core.BaseDir + arguments.EnvYAML.Core.AppName

	// Temp the files to
	arguments.TempDir =  base_dir + arguments.EnvYAML.Core.TempDir
	tmp_bool, err := methods.DoesFileOrDirExists(arguments.TempDir)
	if err != nil {return err}
	if !tmp_bool {
		log.Warn("Directory \""+ arguments.TempDir + "\" does not exists, creating one")
		err:= os.MkdirAll(arguments.TempDir, 0755)
		if err != nil {return err}
	}

	// Download the files to
	arguments.DownloadDir =  base_dir + arguments.EnvYAML.Download.DownloadDir
	dl_bool, err := methods.DoesFileOrDirExists(arguments.DownloadDir)
	if err != nil {return err}
	if !dl_bool {
		log.Warn("Directory \""+ arguments.DownloadDir + "\" does not exists, creating one")
		err:= os.MkdirAll(arguments.DownloadDir, 0755)
		if err != nil {return err}
	}

	// Environment location
	arguments.EnvFileDir = base_dir + arguments.EnvYAML.Install.EnvDir
	env_bool, err := methods.DoesFileOrDirExists(arguments.EnvFileDir)
	if err != nil {return err}
	if !env_bool {
		log.Warn("Directory \""+ arguments.EnvFileDir + "\" does not exists, creating one")
		err := os.MkdirAll(arguments.EnvFileDir, 0755)
		if err != nil {return err}
	}

	// Uninstall location
	arguments.UninstallDir = base_dir + arguments.EnvYAML.Install.UnistallDir
	uninstall_bool, err := methods.DoesFileOrDirExists(arguments.UninstallDir)
	if err != nil {return err}
	if !uninstall_bool {
		log.Warn("Directory \""+ arguments.UninstallDir + "\" does not exists, creating one")
		err := os.MkdirAll(arguments.UninstallDir, 0755)
		if err != nil {return err}
	}

	// Future Reference location
	arguments.FutureRefDir = base_dir + arguments.EnvYAML.Install.FutureRef
	futureref_bool, err := methods.DoesFileOrDirExists(arguments.FutureRefDir)
	if err != nil {return err}
	if !futureref_bool {
		log.Warn("Directory \""+ arguments.FutureRefDir + "\" does not exists, creating one")
		err := os.MkdirAll(arguments.FutureRefDir, 0755)
		if err != nil {return err}
	}

	return nil
}


func Config() error {

	home := os.Getenv("HOME")
	// Read the config file and store the value on a struct
	source, err := ioutil.ReadFile(home + "/.config.yml")
	if err != nil {return err}
	yaml.Unmarshal(source, &arguments.EnvYAML)

	// Creating Directory needed for the program
	err = CreateDir()
	if err != nil {return err}

	return nil
}
