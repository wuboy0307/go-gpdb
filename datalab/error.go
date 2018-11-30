package main

var (
	alreadyExists =
`The vagrant name "%[1]s" is already on our config file
1. Sometime vagrant provision failed and we updated the configuration, try running the below command to remove it from our configuration

%[2]s destroy -n %[1]s 

2. If this was removed manually using "vagrant destroy" or you removed the vm manually from virtual box, then try running the below command to remove it from our configuration

%[2]s delete-config -n %[1]s

3. You can also change the name of the hostname while creating using 

%[2]s create -n <new-hostname> ....
`
	apiTokenMissing =
`The API Token is not set, please run the below command to set it

%s update-config -t <token>
`
	goGPDBLocationMissing =
`The go-gpdb file location is not set, please run the below command to set it 

%s update-config -l <vagrant file location>
`
	missingVMInOurConfig =
`Cannot run "%s", since we don't know anything about the VM with the name "%s", maybe its a typo on the hostname if yes then try again with the correct name. You can also check all the provisioned VM's by the program %[3]s using the command "%[3]s list"`
)
