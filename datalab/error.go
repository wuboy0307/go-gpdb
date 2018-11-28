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
	vagrantLocationMissing =
`The Vagrant Location is not set, please run the below command to set it 

%s update-config -l <vagrant file location>
`
)
