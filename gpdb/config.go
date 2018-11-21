package main

import (
	"github.com/jinzhu/configor"
	"os"
)

// Struct that store the configuration file for the program to run
var Config = struct {
	CORE struct {
		APPLICATIONNAME string `yaml:"APPLICATION_NAME"`
		OS              string `yaml:"OS"`
		ARCH            string `yaml:"ARCH"`
		GOBUILD         string `yaml:"GO_BUILD"`
		BASEDIR         string `yaml:"BASE_DIR"`
		TEMPDIR         string `yaml:"TEMP_DIR"`
	} `yaml:"CORE"`
	DOWNLOAD struct {
		APITOKEN    string `yaml:"API_TOKEN"`
		DOWNLOADDIR string `yaml:"DOWNLOAD_DIR"`
	} `yaml:"DOWNLOAD"`
	INSTALL struct {
		ENVDIR               string `yaml:"ENV_DIR"`
		UNINTSALLDIR         string `yaml:"UNINTSALL_DIR"`
		FUTUREREFDIR         string `yaml:"FUTUREREF_DIR"`
		MASTERUSER           string `yaml:"MASTER_USER"`
		MASTERPASS           string `yaml:"MASTER_PASS"`
		GPMONPASS            string `yaml:"GPMON_PASS"`
		MASTERDATADIRECTORY  string `yaml:"MASTER_DATA_DIRECTORY"`
		SEGMENTDATADIRECTORY string `yaml:"SEGMENT_DATA_DIRECTORY"`
		TOTALSEGMENT         int    `yaml:"TOTAL_SEGMENT"`
	} `yaml:"INSTALL"`
}{}

// Read the configuration file and create directory if not exists
// or set default values if values are missing
func validateConfiguration() {

	Debug("Checking if the directories needed for the program exists")
	if IsValueEmpty(Config.CORE.BASEDIR) {
		Warn("BASE_DIR parameter missing in the config file, setting to default")
		Config.CORE.BASEDIR = "/home/gpadmin/"
	}

	// App name
	if IsValueEmpty(Config.CORE.APPLICATIONNAME) {
		Warn("APPLICATION_NAME parameter missing in the config file, setting to default")
		Config.CORE.APPLICATIONNAME = "gpdbinstall"
	}

	// Temp Directory
	if IsValueEmpty(Config.CORE.TEMPDIR) {
		Warn("TEMP_DIR parameter missing in the config file, setting to default")
		Config.CORE.TEMPDIR = "/temp/"
	}

	// Download Directory
	if IsValueEmpty(Config.DOWNLOAD.DOWNLOADDIR) {
		Warn("DOWNLOAD_DIR parameter missing in the config file, setting to default")
		Config.DOWNLOAD.DOWNLOADDIR = "/download/"
	}

	// Env Directory
	if IsValueEmpty(Config.INSTALL.ENVDIR) {
		Warn("ENV_DIR parameter missing in the config file, setting to default")
		Config.INSTALL.ENVDIR = "/env/"
	}

	// Uninstall Directory
	if IsValueEmpty(Config.INSTALL.UNINTSALLDIR) {
		Warn("UNINTSALL_DIR parameter missing in the config file, setting to default")
		Config.INSTALL.UNINTSALLDIR = "/uninstall/"
	}

	// Future Reference Directory
	if IsValueEmpty(Config.INSTALL.FUTUREREFDIR) {
		Warn("UNINTSALL_DIR parameter missing in the config file, setting to default")
		Config.INSTALL.FUTUREREFDIR = "/future_reference/"
	}

	// Check if API Token
	if IsValueEmpty(Config.DOWNLOAD.APITOKEN) || Config.DOWNLOAD.APITOKEN == "<API TOKEN>" {
		// Check if its set as environment variables
		token := os.Getenv("UAA_API_TOKEN")
		if token == "" { // No token set
			Fatal("The API Token is either missing or not provided, please update the config file and try again")
		} else {
			Config.DOWNLOAD.APITOKEN = token
		}
	}

	// TODO: Check if the MASTER PASS / SEGMENT PASS / MASTER DATA DIRECTORY
	// TODO: TRIM the above stuff

	// Check if the directory exists, else create one.
	base_dir := Config.CORE.BASEDIR + Config.CORE.APPLICATIONNAME

	// Temp the files to
	Config.CORE.TEMPDIR =  base_dir + Config.CORE.TEMPDIR
	CreateDir(Config.CORE.TEMPDIR)

	// Download the files to
	Config.DOWNLOAD.DOWNLOADDIR = base_dir + Config.DOWNLOAD.DOWNLOADDIR
	CreateDir(Config.DOWNLOAD.DOWNLOADDIR)

	// Environment location
	Config.INSTALL.ENVDIR = base_dir + Config.INSTALL.ENVDIR
	CreateDir(Config.INSTALL.ENVDIR)

	// Uninstall location
	Config.INSTALL.UNINTSALLDIR = base_dir + Config.INSTALL.UNINTSALLDIR
	CreateDir(Config.INSTALL.UNINTSALLDIR)

	// Future Reference location
	Config.INSTALL.FUTUREREFDIR = base_dir + Config.INSTALL.FUTUREREFDIR
	CreateDir(Config.INSTALL.FUTUREREFDIR)

}


// Load the configuration file to the memory
func config() {

	// Load the configuration
	configFile := "config.yml"
	Debugf("Reading the configuration file and loading to the memory: %s", configFile)
	configor.Load(&Config, configFile)

	// Validate the configuration
	validateConfiguration()
}