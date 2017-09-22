#!/bin/bash

echo "Reloading sysctl"
sysctl -p

echo "Cleaning RPM cache"
sed -i 's/gpgcheck=1/gpgcheck=0/g' /etc/yum.repos.d/*
sudo yum clean all

echo "Installing RPMs"
sudo yum -y install ed unzip tar git strace gdb vim-enhanced wget m4
sudo yum -y install ed unzip tar git strace gdb vim-enhanced wget m4

echo "Creating gpadmin user"
useradd gpadmin

echo "Setting gpadmin password"
echo "gpadmin:changeme" | chpasswd

echo "Setting root password"
echo "root:changeme" | chpasswd

echo "Cloning piv-go-gpdb repository"
git clone https://github.com/ielizaga/piv-go-gpdb.git /tmp/piv-go-gpdb
chown -R gpadmin:gpadmin /tmp/piv-go-gpdb

echo "Updating the API Token"
sed -i "s|<API TOKEN>|$1|g" /tmp/piv-go-gpdb/config.yml

echo "Changing the permission of /usr/local"
chmod 777 /usr/local

echo "Create /data/ directory"
mkdir -p /data
chown gpadmin:gpadmin /data

echo "Running piv-go-gpdb script"
mkdir /home/gpadmin/piv-go-gpdb
cd /home/gpadmin/piv-go-gpdb
cp /tmp/piv-go-gpdb/config.yml /home/gpadmin/piv-go-gpdb
cp /tmp/piv-go-gpdb/setup.sh /home/gpadmin/piv-go-gpdb
chown -R gpadmin:gpadmin /home/gpadmin
/bin/sh /home/gpadmin/piv-go-gpdb/setup.sh

echo "Moving the go binaries to gpadmin user home directory"
mv $HOME/.go /home/gpadmin
chown gpadmin:gpadmin /home/gpadmin/.go

echo "Copy the bashrc to gpadmin user folder"
cp ~/.bashrc /home/gpadmin

echo "Moving the config file to /home/gpadmin/ directory"
mv ~/.config.yml /home/gpadmin

echo "Updating the vagrant user bashrc to auto login to gpadmin"
{
    echo "sudo su - gpadmin"
    echo "exit"
} >> /home/vagrant/.bashrc

echo "Cleaning up the tmp directory"
rm -rf /tmp/piv-go-gpdb

echo "The vagrant setup is complete"
