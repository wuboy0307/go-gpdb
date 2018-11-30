Table of Contents
=================

   * [Introduction](#introduction)
   * [Creating Vagrant VM's Manually](#creating-vagrant-vms-manually)
   * [Environment](#environment)
        * [PivNet Token](#pivnet-token)
        * [VM Customization](#vm-customization)
        * [Application Defaults](#application-defaults)
        * [Developer mode](#developer-mode)
   * [VM's Setup Files](#vms-setup-files)
        * [functions.h](#functionsh)
        * [os.prep.sh](#osprepsh)
        * [go.build.sh](#gobuildsh)
   * [Developer Mode VM](#developer-mode-vm)

# Introduction

This folder contains a set of bash script that helps to make sure the OS contains basic packages and software that can be used with gpdb softwares.

# Creating Vagrant VM's Manually

If you wish to create you vagrant VM's Manually and don't want use the [datalab CLI](https://github.com/pivotal-gss/go-gpdb/blob/master/datalab/README.md), here are few steps.
+ The vagrant file is located [here](https://github.com/pivotal-gss/go-gpdb/blob/master/Vagrantfile)
+ The vagrant file take in few environment variables more information of the environment variable is described at the [environment section](https://github.com/pivotal-gss/go-gpdb/blob/master/scripts/README.md#environment)
+ During the provision of VM, it runs the two shell script, more information on the shell script is available [here](https://github.com/pivotal-gss/go-gpdb/blob/master/scripts/README.md##vms-setup-files))

# Environment

The [Vagrant file](https://github.com/pivotal-gss/go-gpdb/blob/master/Vagrantfile) can be customized to your liking or based on your needs. If the values are not provided then it chooses it default, the different environment variables includes

### PivNet Token

+ The token to download the gpdb software from pivotal network page
    ```
    # EVN PIVNET DEFAULTS
    @pivnet_token = ENV['UAA_API_TOKEN'] || ""
    ```
    To obtain the pivnet token, Navigate to [PivNet Edit Profile](https://network.pivotal.io/users/dashboard/edit-profile) and scroll to the bottom of the page near “UAA API TOKEN” & click on the button “Request New API Token”, copy the token (**PLEASE NOTE:** This token will change if you click on the “Request New API Token” again)

+ For eg.
    ```
    UAA_API_TOKEN=abcdefgh vagrant up
    ```

### VM Customization 

+ If you are short of RAM or CPU, you can configure using "VM_CPUS", "VM_MEMORY" etc environment variables
    ```
    # EVN VM DEFAULTS
    @vm_os = ENV['VM_OS'] || "bento/centos-7.5"
    @vm_cpus = ENV['VM_CPUS'] || 2
    @vm_memory = ENV['VM_MEMORY'] || 4096
    ```
    **NOTE:** Memory is calculated in MB
+ For eg.
    ```
    VM_OS=xxx VM_CPUS=x VM_MEMORY=x vagrant up
    ```

### Application Defaults
+ You can also customize the name of the VM's, how many VM's do you need ? do you need a standby host ? etc using the below environment variables
    ```
    # ENV APPLICATION DEFAULTS
    @subnet = ENV['GO_GPDB_SUBNET'] || "192.168.99.100"
    @hostname = ENV['GO_GPDB_HOSTNAME'] || "gpdb"
    @segments = ENV['GO_GPDB_SEGMENTS'].to_i || 0
    @standby = ENV['GO_GPDB_STANDBY'] || false
    ```
+ For eg
    ```
    GO_GPDB_SUBNET=192.xxx.xx.xx GO_GPDB_HOSTNAME=xxx/xxx GO_GPDB_SEGMENTS=x GO_GPDB_STANDBY=true vagrant up
    ```

### Developer mode
+ If you want to download go binaries / dep etc then enable the developer mode when provision'ng the VM.
    ```
    # Developer options
    @developer_mode = ENV['DEVELOPER_MODE'] || false
    @go_path = ENV['GOPATH'] ||""
    ```
+ For eg
    ```
    DEVELOPER_MODE=true GOPATH="xxx/xxx" vagrant up
    ```
    
# VM's Setup Files

Here are the details of the files found on the [script](https://github.com/pivotal-gss/go-gpdb/tree/master/scripts) directory

### functions.h

Standard file that has some key function to help the scripts like [os preparation](https://github.com/pivotal-gss/go-gpdb/blob/master/scripts/os.prep.sh) / [go build](https://github.com/pivotal-gss/go-gpdb/blob/master/scripts/go.build.sh) & [test](https://github.com/pivotal-gss/go-gpdb/tree/master/test) to work perfectly

### os.prep.sh

Install basic packages / setup users & set proper permission for the default user to use the provisioned VM

### go.build.sh

Install the go binaries & basic dev packages if in developers mode and also download and sets up the gpdb cli so that the user can start using the VM without any extra work

# Developer Mode VM

To understand more on the developer mode VM, refer to the [README](https://github.com/pivotal-gss/go-gpdb/blob/master/test/README.md) of the test