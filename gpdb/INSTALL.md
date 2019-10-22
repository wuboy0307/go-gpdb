+ Download the latest build of go-gpdb from the [release link](https://github.com/pivotal-gss/go-gpdb/releases/tag/v3.3.0)
+ Create a `config.yml` file at your home directory
```
touch ~/config.yml
```
+ And copy a sample config shown [here](https://github.com/pivotal-gss/go-gpdb/blob/master/gpdb/config.yml), you may change the location of the directory if you would like to.
+ Prepare your OS
  + Make sure `ssh $HOSTNAME` works. If not, you could run `ssh-keygen -f ~/.ssh/id_rsa -N '' && cat .ssh/id_rsa.pub >> .ssh/authorized_keys`
  + Make sure you are logged in as a non-root user ( like gpadmin ) which has sudo privileges (necessary for yum installing the greenplum-db .rpm)
  + Make sure that the non-root user has permission to read / write / execute on the directories listed at the file `config.yml`
  + You can checkout [this script](https://github.com/pivotal-gss/go-gpdb/blob/master/scripts/os.prep.sh) for all the packages and permission needed on the host or to setup your host for gpdb installation.
  + The first line of the `/etc/hosts` is always ignored and it can be anything, your hostname should start from the second line and its usually where the localhost information should starts
  
    **NOTE:** Comments on `/etc/hosts` are not ignored by the tool, its a drawback so ensure there is no comments on /etc/hosts after the first line
    
    Your `/etc/hosts` would have to looks something like this.
    ```
    [gpadmin@gpdb-m ~]$ cat /etc/hosts
    127.0.0.1 localhost 
    192.168.99.100 gpdb-m
    ```
    
    If you have multiple hosts and want to create a cluster, then add them like this below
    ```
    [gpadmin@gpdb-m ~]$ cat /etc/hosts
    127.0.0.1 localhost 
    192.168.99.100 gpdb-m
    192.168.99.101 gpdb1-m
    192.168.99.102 gpdb2-m
    ```
    The tool will detect that there is multiple host and create a cluster, no extra steps needed
        
+ Connect to your non-root user & Run commands
  ```
  # Be at your home directory
  cd
  
  # Install the sshpass file
  sudo yum install -y sshpass
  
  # Set your hostname
  echo $HOSTNAME > hostfile
  
  # If you have a different OS username 
  sed -i "s/gpadmin/$USER/g" config.yml
  
  # Change permission on the binary
  chmod +x gpdb
  
  # Open source gpdb release
  ./gpdb download -v 6.0.1 --github --install 
  
  # Official gpdb release 
  ./gpdb download -v 6.0.1 --install 
  ```

Use your new gpdb cluster! 😉🥳