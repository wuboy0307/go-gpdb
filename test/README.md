Table of Contents
=================

   * [Introduction](#introduction)
   * [Setup](#setup)
   * [Understanding the output](#understanding-the-output)

# Introduction

After the modification is done, it always better to test the basic functionality works, this simple bash script just does that, it connects to pivnet and collects the gpdb software and try to install and verify all the version installation works without any issue.

# Setup 

+ Add the below line on the Vagrantfile, to mount the local drive to the VM that is provisioned by vagrant

```
# All Vagrant configuration is done below. 
Vagrant.configure("2") do |config|

  config.vm.synced_folder "/Users/xxxx/Documents/Project/go-gpdb/", "/gpdb"

  @ip = IPAddr.new @subnet

  config.vm.box = @vm_os
```

+ Now create a brand new VM with the below configuration

```
datalab create -s 2 --standby
```

for more information on the datalab cli refer to the documentation here

+ Once provision connect and navigate to the directory 

```
datalab ssh
cd /gpdb/src/github.com/pivotal-gss/go-gpdb/gpdb
```

+ Setup the path 

```
PATH=$PATH:$HOME/.local/bin:$HOME/bin:/usr/local/go/bin
export GOROOT=/usr/local/go
export GOPATH=/gpdb/
export PATH
```

+ And download & install binaries

```
go run *.go download -v 5.13.0
go run *.go install -v 5.13.0
```

The reason why we need this step is to ensure we have a password less access so that the test can run without our intervention

+ Now remove the installation 

```
go run *.go remove -v 5.13.0
```

+ And now go to the test directory

```
cd /gpdb/src/github.com/pivotal-gss/go-gpdb/test
```

+ Now run the test and wait for it to succeed in all the steps

```
/bin/bash gpdb.unit.test.sh
```

# Understanding the output

A line with a tick [√] indicates all the checks passed for the version for eg.s

```
[gpadmin@gpdb-m test]$ /bin/bash gpdb.unit.test.sh

============================================================
TEST: GO GPDB Version: "5.13.0"
============================================================
[√] Downloading the version 5.13.0...
[√] Install the version 5.13.0...
[√] Checking if the database with version 5.13.0 is healthy...
[√] Downloading the cc version 4.4.2 for gpdb version 5.13.0...
[√] Install the cc version 4.4.2 on gpdb Version 5.13.0...
[√] Checking if the CC Url 4.4.2 on gpdb Version 5.13.0 is working...
[√] Setting the environment of version 5.13.0...
[√] Remove the version 5.13.0...
[√] Cleaning the binaries in the /usr/local & /usr/local/src/gpdbinstall/download/...
[.....]
```

a cross [X] indicated that checks didn't pass 

for eg.s

```
[X] Checking if the CC Url 2.0 on gpdb Version 4.3.1.0 is working...
[X] Return Code: [1]...
[√] Downloading the cc version 2.1.0 for gpdb version 4.3.1.0...
[X] Install the cc version 2.1.0 on gpdb Version 4.3.1.0...
[X] Return Code: [1]...
[X] Checking if the CC Url 2.1.0 on gpdb Version 4.3.1.0 is working...
[X] Return Code: [1]...
```

and you will need to check the logs to understand the reason, the logs are available on the location "/tmp"