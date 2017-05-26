#!/bin/bash

echo "Reloading sysctl"
sysctl -p

echo "Cleaning RPM cache"
sed -i 's/gpgcheck=1/gpgcheck=0/g' /etc/yum.repos.d/*
sudo yum clean all

echo "Installing RPMs"
sudo yum -y install ed unzip tar git
sudo yum -y install ed unzip tar git

echo "Creating gpadmin user"
useradd gpadmin

echo "Setting gpadmin password"
echo "gpadmin:changeme" | chpasswd

echo "Cloning piv-go-gpdb repository"
git clone https://github.com/ielizaga/piv-go-gpdb.git /home/gpadmin/piv-go-gpdb
chown -R gpadmin:gpadmin /home/gpadmin

echo "Updating the API Token"
sed -i "s|<API TOKEN>|$1|g" /home/gpadmin/piv-go-gpdb/config.yml

echo "Changing the permission of /usr/local"
chmod 777 /usr/local

echo "Create /data/ directory"
mkdir -p /data
chown gpadmin:gpadmin /data

echo "Running piv-go-gpdb script"
cd /home/gpadmin/piv-go-gpdb
/bin/sh setup.sh

echo "Moving the go binaries to gpadmin user home directory"
mv $HOME/.go /home/gpadmin
chown gpadmin:gpadmin /home/gpadmin/.go

echo "Copy the bashrc to gpadmin user folder"
cp ~/.bashrc /home/gpadmin

echo "Moving the config file to /home/gpadmin/ directory"
mv ~/.config.yml /home/gpadmin

echo "Updating the vagrant user bashrc to auto login to gpadmin"
echo "sudo su - gpadmin" >> /home/vagrant/.bashrc

echo "The vagrant setup is complete"
