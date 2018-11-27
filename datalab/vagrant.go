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
		Fatalf("The vagrant name \"%[1]s\" is already on our config file \n" +
			"1. Sometime vagrant provision failed and we updated the configuration, try running \"%[2]s delete-config -n %[1]s\" to remove it from our configuration \n" +
			"2. If this was removed using \"vagrant destroy\" or you removed the vm manually from virtual box, " +
			"then try running \"%[2]s delete-config -n %[1]s\" to remove it from our configuration \n" +
			"3. You can also change the name of the hostname while creating using \"%[2]s create -n <new-hostname> ....\"", cmdOptions.Hostname, programName)
	}
	executeOsCommand("vagrant", generateEnvArray(), "up")
	Config.Vagrants = append(Config.Vagrants, vk)
}
