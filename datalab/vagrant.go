package main

import "fmt"

// Generate the environment
func generateEnvArray() []string {
	return []string{
		fmt.Sprintf("VM_OS=%s", cmdOptions.Os),
		fmt.Sprintf("VM_CPUS=%d", cmdOptions.Cpu),
		fmt.Sprintf("VM_MEMORY=%d", cmdOptions.Memory),
		fmt.Sprintf("GO_GPDB_SUBNET=%s", cmdOptions.Subnet),
		fmt.Sprintf("GO_GPDB_HOSTNAME=%s", cmdOptions.Hostname),
		fmt.Sprintf("GO_GPDB_SEGMENTS=%d", cmdOptions.Segments),
		fmt.Sprintf("GO_GPDB_STANDBY=%t", cmdOptions.Standby),
	}
}

// Create the VM
func createVM() {
	vk := VagrantKey{
		cmdOptions.Hostname,
		cmdOptions.Cpu,
		cmdOptions.Memory,
		cmdOptions.Standby,
		cmdOptions.Os,
		cmdOptions.Subnet,
		cmdOptions.Segments,
	}
	if nameInConfig(cmdOptions.Hostname) {
		Fatalf(alreadyExists, cmdOptions.Hostname, programName)
	}
	executeOsCommand("vagrant", generateEnvArray(), "up")
	Config.Vagrants = append(Config.Vagrants, vk)
}
