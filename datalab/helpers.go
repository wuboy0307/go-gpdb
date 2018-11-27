package main

import (
	"strings"
	"os/exec"
	"os"
)

// Check is the value is empty
func IsValueEmpty(v string) bool {
	if len(strings.TrimSpace(v)) == 0 {
		return true
	}
	return false
}

// Check if the Os executable exists
func isCommandAvailable(name string) {
	cmd := exec.Command("command", name, "-v")
	if err := cmd.Run(); err != nil {
		Fatalf("the command executable \"%s\" is not installed on your machine, please install it", name)
	}
}

// Check if we have the name already on our config
func nameInConfig(a string) bool {
	for _, b := range Config.Vagrants {
		if b.Name == a {
			return true
		}
	}
	return false
}

// Execute the OS command
func executeOsCommand(command string, env []string, args ...string) {
	// Check if the command is installed on the local host
	isCommandAvailable(command)
	// Initialize the command
	cmd := exec.Command(command, args...)
	// Set all the env variable
	cmd.Env = os.Environ()
	for _, e := range env {
		cmd.Env = append(cmd.Env, e)
	}
	// Attach the stdout and std err
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	// Start and wait for the execution
	err := cmd.Start()
	if err != nil {
		Fatalf("Failed to start the command %s, with arguments %s, err: %v", command, args, err)
	}
	err = cmd.Wait()
	if err != nil {
		Fatalf("Failed in waiting for the command %s, with arguments %s, err: %v", command, args, err)
	}
}