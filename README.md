# Introduction [![Go Version](https://img.shields.io/badge/go-v1.11.1-green.svg?style=flat-square)](https://golang.org/dl/) [![MIT License](https://img.shields.io/badge/License-MIT-red.svg?style=flat-square)](https://github.com/ielizaga/piv-go-gpdb/blob/master/LICENSE)

This reposistory is split into two parts

+ Gpdb cli 

    The gpdb cli helps to download, install, remove and manage the software of GPDB / GPCC.

+ Datalab cli
    
    The datalab cli helps to create & manage vagrant VM's provisioning.
    
Table of Contents
=================

   * [Prerequisite](#prerequisite)
        * [VirtualBox](#virtualbox)
        * [Vagrant](#vagrant)
   * [Tools](#tools)
        * [Gpdb CLI](#gpdb-cli)
        * [Datalab CLI](#datalab-cli)
        * [UnitTest](#unittest)
        * [Vagrant](#vagrant-1)
   * [Developers / Contributor's](#developers--contributors) 

# Prerequisite

The "go-gpdb" software needs the below two tools pre-installed on your machine for it to work.

+ Vagrant
+ VirtualBox

Please follow the below instruction on how to setup the prerequisite

### VirtualBox

+ If not already installed, [download](http://download.virtualbox.org/virtualbox/5.1.22/VirtualBox-5.1.22-115126-OSX.dmg) and install VirtualBox or you can use brew to install virtual box using the command

```
brew update
brew cask install virtualbox
```

+ Once you virtualbox installation is complete ensure you have two interfaces (namely vboxnet0/1 is seen) on your MAC,

```
IRFALI123:Vagrant fali$ ifconfig
vboxnet0: flags=8842<BROADCAST,RUNNING,SIMPLEX,MULTICAST> mtu 1500
	ether 0a:00:27:00:00:00
vboxnet1: flags=8943<UP,BROADCAST,RUNNING,PROMISC,SIMPLEX,MULTICAST> mtu 1500
	ether 0a:00:27:00:00:01
	inet 192.168.11.1 netmask 0xffffff00 broadcast 192.168.11.255
```

if its not shown then refer to the [link](http://islandora.ca/content/fixing-missing-vboxnet0) for Virtualbox version lower than 5 and if your virtualbox version is 5 and above follow this [link](https://luppeng.wordpress.com/2017/07/17/enabling-virtualbox-host-only-adapter-on-mac-os-x/) on how to set those two interfaces up

### Vagrant

+ On your MAC install vagrant using the below command ( if vagrant executable is not already installed )

```
brew update
brew cask install vagrant
```

If you have already installed vagrant ensure you are running the latest version of vagrant, to update your vagrant run

```
brew update
brew cask reinstall vagrant
```

+ Once the vagrant is installed, install the vagrant plugin using the below command

```
vagrant plugin install vagrant-hosts
```

# Tools
 
### Gpdb cli

+ Please check the gpdb cli README for details on how to install & use the gpdb cli.

### Datalab cli

+ Please check the datalab cli README for details on how to install & use the datalab cli.

### UnitTest

+ Please check the README on how to run the unit test case.

### Vagrant 

+ If you wish to install vagrant manually using the Vagrant file & don't want to use the datalab cli, please follow the instruction mentioned here for all the options.

# Developers / Contributor's

Please read the section on how to setup the environment to test and hack this tool