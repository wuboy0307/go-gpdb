package main

import "fmt"

// Generate the environment
func generateEnvArray(cpu, memory, segments int, os, subnet, hostname string, standby, developer bool) []string {
	return []string{
		fmt.Sprintf("UAA_API_TOKEN=%s", Config.APIToken),
		fmt.Sprintf("VM_OS=%s", os),
		fmt.Sprintf("VM_CPUS=%d", cpu),
		fmt.Sprintf("VM_MEMORY=%d", memory),
		fmt.Sprintf("GO_GPDB_SUBNET=%s", subnet),
		fmt.Sprintf("GO_GPDB_HOSTNAME=%s", hostname),
		fmt.Sprintf("GO_GPDB_SEGMENTS=%d", segments),
		fmt.Sprintf("GO_GPDB_STANDBY=%t", standby),
		fmt.Sprintf("DEVELOPER_MODE=%t", developer),
	}
}

// Create the VM
func createVM() {
	Debugf("Provisioning of creating a new vagrant VM with the name: %s", cmdOptions.Hostname)
	vk := VagrantKey{
		cmdOptions.Hostname,
		cmdOptions.Cpu,
		cmdOptions.Memory,
		cmdOptions.Standby,
		cmdOptions.Os,
		cmdOptions.Subnet,
		cmdOptions.Segments,
		cmdOptions.Developer,
	}
	_, exists := nameInConfig(cmdOptions.Hostname)
	if exists {
		Fatalf(alreadyExists, cmdOptions.Hostname, programName)
	}
	Config.Vagrants = append(Config.Vagrants, vk)
	env := generateEnvArray(cmdOptions.Cpu, cmdOptions.Memory, cmdOptions.Segments, cmdOptions.Os, cmdOptions.Subnet, cmdOptions.Hostname, cmdOptions.Standby, cmdOptions.Developer)
	executeOsCommand("vagrant", env, "up")
}

// ssh the VM
func sshVM() {
	Debugf("Ssh to vagrant VM with the name: %s", cmdOptions.Hostname)
	_, exists := nameInConfig(cmdOptions.Hostname)
	if !exists {
		Fatalf(missingVMInOurConfig, "ssh", cmdOptions.Hostname, programName)
	}
	sshInfo := []string{
		fmt.Sprintf("VAGRANT_DOTFILE_PATH=%[1]s/.vagrant VAGRANT_VAGRANTFILE=%[1]s/Vagrantfile vagrant ssh %s-m", removeSlash(Config.GoGPDBPath), cmdOptions.Hostname),
	}
	printOnScreen("Copy & Paste the below information on the terminal to connect to the VM", sshInfo)
}

// destroy the VM
func destroyVM() {
	Debugf("Destroying the vagrant VM with the name: %s", cmdOptions.Hostname)
	executeOsCommand("vagrant", generateEnvByVagrantKey("destroy"), "destroy", "-f")
	deleteConfigKey()
}

// status the VM
func statusVM() {
	Debugf("Status of the vagrant VM with the name: %s", cmdOptions.Hostname)
	executeOsCommand("vagrant", generateEnvByVagrantKey("status"), "status")
}

// bring up the VM
func upVM() {
	Debugf("Bringing up the vagrant VM with the name: %s", cmdOptions.Hostname)
	executeOsCommand("vagrant", generateEnvByVagrantKey("up"), "up")
}

// stop the VM
func stopVM() {
	Debugf("Stopping the vagrant VM with the name: %s", cmdOptions.Hostname)
	executeOsCommand("vagrant", generateEnvByVagrantKey("stop"), "suspend")
}

// restart the VM
func restartVM() {
	Debugf("Restart the vagrant VM with the name: %s", cmdOptions.Hostname)
	executeOsCommand("vagrant", generateEnvByVagrantKey("restart"), "reload")
}

// list the vms
func listVM() {
	Debugf("Listing all the vagrant VM which has been provisioned")
	// Generate the report of the VM's we have recorded
	var output = []string{`Index | Name | CPU | Memory | StandbyVM | OS | Subnet | Segments`,
		`---------|-----------------|---------| ---------|---------|----------------------|--------------------|---------|`}
	for i, vagrant := range Config.Vagrants {
		output = append(output, fmt.Sprintf("%d|%s|%d|%d MB|%t|%s|%s|%d", i+1, vagrant.Name, vagrant.CPU, vagrant.Memory, vagrant.Standby, vagrant.Os, vagrant.Subnet, vagrant.Segment))
	}
	printOnScreen("Here is the list of all the Provisioned Vagrant VM's we know about.", output)

	// If requested to view all the global status then print the vagrant global status
	if cmdOptions.GlobalStatus {
		printOnScreen("Here is the global status information from the vagrant CLI", []string{})
		globalStatus()
	}
}

// print the global status output using the vagrant CLI
func globalStatus() {
	Debugf("Output the vagrant global status information")
	executeOsCommand("vagrant", []string{}, "global-status")
}