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
	VagrantFile string `json:"vagrant_file"`
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
	// If user asked to update the toke
	if !IsValueEmpty(cmdOptions.Token) {
		Config.APIToken = cmdOptions.Token
	}
	// If user asks to update the vagrant location
	if !IsValueEmpty(cmdOptions.Location) {
		Config.VagrantFile = cmdOptions.Location
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

// Load all the configuration
func config() {
	createConfig() // will only create if not found
	err := configor.Load(&Config, configFile)
	if err != nil {
		Fatalf("Failed to load the configuration file %s, err: %v", configFile, err)
	}
}