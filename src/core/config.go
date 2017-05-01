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

	// Uninstall Directory
	if methods.IsValueEmpty(arguments.EnvYAML.Install.UnistallDir) {
		log.Warn("UNINTSALL_DIR parameter missing in the config file, setting to default")
		arguments.EnvYAML.Install.UnistallDir = "/uninstall/"
	}

	// Check if the directory exists, else create one.
	base_dir := arguments.EnvYAML.Core.BaseDir + arguments.EnvYAML.Core.AppName

	// Download the files to
	arguments.DownloadDir =  base_dir + arguments.EnvYAML.Download.DownloadDir
	dl_bool, err := methods.DoesFileOrDirExists(arguments.DownloadDir)
	methods.Fatal_handler(err)
	if !dl_bool {
		log.Warn("Directory \""+ arguments.DownloadDir + "\" does not exists, creating one")
		err:= os.MkdirAll(arguments.DownloadDir, 0755)
		methods.Fatal_handler(err)
	}

	// Environment location
	arguments.EnvFileDir = base_dir + arguments.EnvYAML.Install.EnvDir
	env_bool, err := methods.DoesFileOrDirExists(arguments.EnvFileDir)
	methods.Fatal_handler(err)
	if !env_bool {
		log.Warn("Directory \""+ arguments.EnvFileDir + "\" does not exists, creating one")
		err := os.MkdirAll(arguments.EnvFileDir, 0755)
		methods.Fatal_handler(err)
	}

	// Environment location
	arguments.UnistallDir = base_dir + arguments.EnvYAML.Install.UnistallDir
	uninstall_bool, err := methods.DoesFileOrDirExists(arguments.UnistallDir)
	methods.Fatal_handler(err)
	if !uninstall_bool {
		log.Warn("Directory \""+ arguments.UnistallDir + "\" does not exists, creating one")
		err := os.MkdirAll(arguments.UnistallDir, 0755)
		methods.Fatal_handler(err)
	}
}


func Config() {

	// Read the config file and store the value on a struct
	source, err := ioutil.ReadFile("config.yml")
	methods.Fatal_handler(err)
	yaml.Unmarshal(source, &arguments.EnvYAML)

	// Creating Directory needed for the program
	CreateDir()
}
