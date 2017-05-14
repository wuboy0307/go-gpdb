#!/usr/bin/env bash
#!/bin/bash
set -e

# things to do after install
# Yum update
# yum install git -y
# clone the git: https://github.com/ielizaga/piv-go-gpdb.git

#
# Core check
#

# Check if the internet connect is working
wget -q --tries=10 --timeout=20 --spider http://google.com
if [[ $? -eq 0 ]]; then
        echo "Internet connection is available, contine ...."
else
        echo "No Internet connection, exiting from the program..."
        exit 2
fi

#
# Download and install GO Binaries.
#

# Setting up go version to download
VERSION="1.7.4"
DFILE="go$VERSION.linux-amd64.tar.gz"

# If the version of go already exit then uninstall it
if [ -d "$HOME/.go" ]; then
        rm -rf $HOME/.go
fi

# Downloading the go tar file
echo "Downloading $DFILE ..."
wget https://storage.googleapis.com/golang/$DFILE -O /tmp/go.tar.gz
if [ $? -ne 0 ]; then
    echo "Download failed! Exiting."
    exit 1
fi

# Extracting the file
echo "Extracting ..."
tar -C "$HOME" -xzf /tmp/go.tar.gz
mv "$HOME/go" "$HOME/.go"

# Updating the bashrc with the information of GOROOT.
if grep -q "GOROOT" "$HOME/.bashrc";
then
    echo "GOROOT binaries location is already updated on the .bashrc file"
else
    touch "$HOME/.bashrc"
    {
        echo '# Golang binaries'
        echo 'export GOROOT=$HOME/.go'
        echo 'export PATH=$PATH:$GOROOT/bin'
    } >> "$HOME/.bashrc"
fi

# Update bashrc with the information of GOPATH.
if grep -q "GOPATH" "$HOME/.bashrc";
then
    echo "GOPATH location is already updated on the .bashrc file"
else
    pwd=`pwd`
    touch "$HOME/.bashrc"
    {
        echo '# GOPATH location'
        echo 'export GOPATH='${pwd}
        echo 'export PATH=$PATH:$GOPATH/bin'
    } >> "$HOME/.bashrc"
fi

# Remove the downloaded tar file
rm -f /tmp/go.tar.gz

#
# Upgrading the code (if any)
#

echo "Pulling newer version of the code"
cp config.yml /tmp/config.yml
git clean -fd
git pull
mv /tmp/config.yml config.yml


#
# Download program dependencies
#

echo "Downloading program dependencies"

# YAML package
source "$HOME/.bashrc"
go get gopkg.in/yaml.v2
if [ $? -ne 0 ]; then
    echo "Download failed, if dependencies package failed. Exiting....."
    exit 1
fi

#
# Build go executable file.
#

echo "Compiling the program... "

# Compile the program
go build gpdb.go
if [ $? -ne 0 ]; then
    echo "Cannot build gpdb executable, exiting ....."
    exit 1
fi

# move the binary file to bin directory
if [ ! -d bin ]; then
    mkdir -p bin/
fi

# move it to bin directory
mv gpdb bin/

#
# Success message.
#

echo "GPDBInstall Script has been successfully installed"
echo "Please close this terminal and open up a new terminal to set the environment"
