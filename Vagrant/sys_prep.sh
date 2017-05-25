#!/bin/bash

echo "Reloading sysctl"
sysctl -p

echo "Cleaning RPM cache"
yum clean all

echo "Installing RPMs"
yum -y install ed unzip tar git

echo "Creating gpadmin user"
useradd gpadmin

echo "Setting gpadmin password"
echo "gpadmin:changeme" | chpasswd

echo "Cloning piv-go-gpdb repository"

git clone https://github.com/ielizaga/piv-go-gpdb.git /home/gpadmin/piv-go-gpdb
chown -R gpadmin:gpadmin /home/gpadmin

sed -i "s|<API TOKEN>|$1|g" /home/gpadmin/piv-go-gpdb/config.yml
