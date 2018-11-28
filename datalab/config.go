package main

import (
	"github.com/jinzhu/configor"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
)

type ProgramConfig struct {
	APIToken    string `json:"api_token"`
	GoGPDBPath string `json:"go_gpdb_path"`
	Vagrants    []VagrantKey `json:"vagrants,omitempty"`
}

type VagrantKey struct {
	Name     string `json:"name,omitempty"`
	CPU      int    `json:"cpu,omitempty"`
	Memory   int    `json:"memory,omitempty"`
	Standby  bool   `json:"standby"`
	Os       string `json:"os,omitempty"`
	Subnet   string `json:"subnet,omitempty"`
	Segment  int    `json:"segment"`
}

var (
	Config ProgramConfig
	configFile = fmt.Sprintf("%s/%s", os.Getenv("HOME"), configFileName)
)

// Create the configuration file if it doesn't exits
func createConfig() {
	Debugf("Checking & Creating the file %s, if not exists", configFile)
	// detect if file exists
	var _, err = os.Stat(configFile)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(configFile)
		if err != nil {
			Fatalf("Error in creating a file: %v", err)
		}
		defer file.Close()
	}
}

// Save Configuration
func saveConfig() {
	Debugf("Saving the configuration information in the config file: %s", configFile)
	// First lets clear all the old configuration
	deleteConfig()
	// Create a new file
	createConfig()
	// Now lets start to write the configuration file
	configJson, err := json.MarshalIndent(Config, "", "   ")
	if err != nil {
		Fatalf("Error when json marshal of the config file, err: %v", err)
	}
	err = ioutil.WriteFile(configFile, configJson, 0644)
	if err != nil {
		Fatalf("Error when creating the config file %s, err: %v", configFile, err)
	}
}

// Update Configuration
func updateConfig() {
	Debugf("Updating the configuration information to be updated config file: %s", configFile)
	// If user asked to update the toke
	if !IsValueEmpty(cmdOptions.Token) {
		Config.APIToken = cmdOptions.Token
	}
	// If user asks to update the go-gpdb location
	if !IsValueEmpty(cmdOptions.GoGPDBPath) {
		Config.GoGPDBPath = cmdOptions.GoGPDBPath
	}
}

// Delete Configuration
func deleteConfig() {
	Debugf("Deleting the file: %s", configFile)
	// delete file
	var err = os.RemoveAll(configFile)
	if err != nil && !os.IsNotExist(err) {
		Fatalf("Error in deleting the file: %v", err)
	}
}

// Delete Configuration key is vm destroy'ed or if requested
func deleteConfigKey() {
	Debugf("Deleting the configuration with the VM name: %s", cmdOptions.Hostname)
	index, exists := nameInConfig(cmdOptions.Hostname)
	if !exists {
		Fatalf(missingVMInOurConfig, "delete-config", cmdOptions.Hostname, programName)
	}
	Config.Vagrants = append(Config.Vagrants[:index], Config.Vagrants[index+1:]...)
	saveConfig()
}

// Load all the configuration
func config() {
	Debugf("Loading all the configuration information")
	createConfig() // will only create if not found
	err := configor.Load(&Config, configFile)
	if err != nil {
		Fatalf("Failed to load the configuration file %s, err: %v", configFile, err)
	}
}