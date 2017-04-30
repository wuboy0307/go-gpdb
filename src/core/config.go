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

func CreateDir() {

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

	// Check if the directory exists, else create one.
	base_dir := arguments.EnvYAML.Core.BaseDir + arguments.EnvYAML.Core.AppName

	// Download the files to
	arguments.DownloadDir =  base_dir + arguments.EnvYAML.Download.DownloadDir
	dl_bool, _ := methods.DoesDirexists(arguments.DownloadDir)
	if !dl_bool {
		log.Warn("Directory \""+ arguments.DownloadDir + "\" does not exists, creating one")
		os.MkdirAll(arguments.DownloadDir, 0777)
	}

	// Environment location
	arguments.EnvFileDir = base_dir + arguments.EnvYAML.Install.EnvDir
	env_bool, _ := methods.DoesDirexists(arguments.EnvFileDir)
	if !env_bool {
		log.Warn("Directory \""+ arguments.EnvFileDir + "\" does not exists, creating one")
		os.MkdirAll(arguments.EnvFileDir, 0777)
	}
}


func Config() {

	// Read the config file and store the value on a struct
	source, err := ioutil.ReadFile("config.yml")
	methods.Fatal_handler(err)
	yaml.Unmarshal(source, &arguments.EnvYAML)

	// Creating Directory
	CreateDir()

}
