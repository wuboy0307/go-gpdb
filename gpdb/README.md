Table of Contents
=================

   * [Introduction](#introduction)
   * [Prerequisite](#prerequisite)
   * [Installation](#installation)
   * [Installation Demo](#installation-demo)
   * [Usage](#usage)
   * [Example](#example)
        * [Download](#download)
        * [Install](#install)
        * [Env](#env)
        * [Remove](#remove)
   * [Demo](#demo)
   * [Developers / Contributors](#developers--contributors)

# Introduction

Gpdb CLI helps to download, install, remove and manage the software of GPDB / GPCC, either as single node installation or multi node installation.

# Prerequisite

+ You will need to first provision the vagrant VM, please follow the instruction as described on the [datalab documentation](https://github.com/pivotal-gss/go-gpdb/tree/master/datalab#create) to provision a vagrant VM.
+ Make sure your provisioned can access the internet.

# Installation

+ The gpdb cli is automatically downloaded and installed on the provisioned VM, there is no additional step needed.
+ Connect to the provisioned VM using [datalab ssh](https://github.com/pivotal-gss/go-gpdb/tree/master/datalab#ssh) and start using the gpdb cli using the [examples](#example) mentioned below

# Usage 

The usage of gpdb CLI

```
Usage:
  gpdb [command] [flags]
  gpdb [command]

Available Commands:
  download    Download the product from pivotal network
  env         Show all the environment installed
  help        Help about any command
  install     Install the product downloaded from download command
  remove      Removes the product installed using the install command

Flags:
  -d, --debug     Enable verbose or debug logging
  -h, --help      help for gpdb
      --version   version for gpdb

Use "gpdb [command] --help" for more information about a command.
```

# Example

### Download

+ To download product interactively
    ```
    gpdb download
    ```
+ To download a specific version
    ```
    gpdb download -v <GPDB VERSION>
    ```
+ To download and install a specific version
    ```
    gpdb download -v <GPDB VERSION> --install
    ```
+ To download GPCC software in interactive mode.
    ```
    gpdb download -p gpcc
    ```
+ To download GPCC software of specific version.
    ```
    gpdb download -p gpcc -v <GPDB VERSION>
    ```
+ To download all GPDB products in interactive mode
    ```
    gpdb download -p gpextras
    ```
+ To download all products of specific version.
    ```
    gpdb download -p gpextras -v <GPDB VERSION>
    ```
+ To obtain help menu of the download command
    ```
    gpdb help download
    ```

### Install

+ To install gpdb
    ```
    gpdb install -v <GPDB VERSION>
    ```
+ To install gpdb & standby
    ```
    gpdb install -v <GPDB VERSION> --standby
    ```
+ To install gpcc
    ```
    gpdb install -p gpcc -v <GPDB VERSION> -c <GPCC VERSION>
    ```
+ To obtain help menu of the install command
    ```
    gpdb help install
    ```

### Env

+ To list all environment that has been installed.
    ```
    gpdb env -l
    ```
+ To list all environment that has been installed and choose env in interactive mode.
    ```
    gpdb env
    ```
+  To start and use a specific installation.
    ```
    gpdb env -v <GPDB VERSION>
    ```
+ To prevent stopping other environment when the environment is set.
    ```
    gpdb env -v <GPDB VERSION> --dont-stop
    ```
+ To obtain help menu of the env command
    ```
    gpdb help env
    ```
    
### Remove

+ To remove a particular installation.
    ```
    gpdb remove -v <GPDB VERSION>
    ```
+ To remove a particular installation forcefully.
    ```
    gpdb remove -v <GPDB VERSION> -f
    ```
+ To obtain help menu of the remove command
    ```
    gpdb help remove
    ```

# Demo

[![asciicast](https://asciinema.org/a/HqncgdNd3CmuexNSHbXmtrL4w.svg)](https://asciinema.org/a/HqncgdNd3CmuexNSHbXmtrL4w)

# Developers / Contributors

1. Clone the git repository
2. Export the GOPATH
    ```
    export GOPATH=<path to the clone repository>
    ```
3. cd to project folder
    ```
    cd $GOPATH/src/github.com/pivotal-gss/go-gpdb/gpdb
    ```
4. Install all the dependencies. If you don't have dep installed, follow the instruction from [here](https://github.com/golang/dep)
    ```
    dep ensure
    ```
5. You are all set, you can run it locally using
    ```
    go run *.go <commands>
    ```
6. To build the package use
    ```
    env GOOS=linux GOARCH=amd64 go build -o gpdb
    ```