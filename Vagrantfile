# -*- mode: ruby -*-
# vi: set ft=ruby :

# Includes
require "ipaddr"

# Developer options
@developer_mode = ENV['DEVELOPER_MODE'] || false
@go_path = ENV['GOPATH'] ||""

# EVN PIVNET DEFAULTS
@pivnet_token = ENV['UAA_API_TOKEN'] || ""

# EVN VM DEFAULTS
@vm_os = ENV['VM_OS'] || "bento/centos-7.5"
@vm_cpus = ENV['VM_CPUS'] || 2
@vm_memory = ENV['VM_MEMORY'] || 4096

# ENV APPLICATION DEFAULTS
@subnet = ENV['GO_GPDB_SUBNET'] || "192.168.99.100"
@hostname = ENV['GO_GPDB_HOSTNAME'] || "gpdb"
@segments = ENV['GO_GPDB_SEGMENTS'].to_i || 0
@standby = ENV['GO_GPDB_STANDBY'] || false

# Define a Template for Building All Our VMs.
def build_vm( config, hostname, ip )
    config.vm.define hostname do |node|
      node.vm.hostname = hostname
      node.vm.network :private_network, :ip => ip

      node.vm.provision :hosts do |provisioner|
         provisioner.autoconfigure = true
         provisioner.sync_hosts = true
         provisioner.add_localhost_hostnames = false
       end

      node.vm.provider :virtualbox do |vb|
          vb.name = hostname
          vb.gui = false
          vb.cpus = @vm_cpus
          vb.memory = @vm_memory
      end
    end
end

# All Vagrant configuration is done below.
Vagrant.configure("2") do |config|

  # Share the project folder if the developer mode is on
  if (@developer_mode == 'true')
    config.vm.synced_folder "#{@go_path}", "/gpdb"
  end

  @ip = IPAddr.new @subnet

  config.vm.box = @vm_os

  # If "vagrant ssh", login As gpadmin, without hacking the vagrant profile
  if ARGV[0] == "ssh"
    config.ssh.username = 'gpadmin'
  end

  # If the token is empty and we are deploying, prompt for a token...
  if (@pivnet_token.to_s.empty?) && (ARGV[0] == 'up')
    puts "","UAA_API_TOKEN Environment Variable Not Found..."
    puts "A NULL entry will deploy, but a token is required for authentication.",""
    print "UAA API TOKEN: "
    @pivnet_token = STDIN.gets.chomp
  end

  puts "","UAA_API_TOKEN: #{@pivnet_token}"
  puts "Master Hostname: #{@hostname}-m"
  puts "Master IP: #{@ip}",""

  # Master Node:
  build_vm( config, "#{@hostname}-m", "#{@ip}" )

  # Create standby host is asked
  if (@standby == 'true')
    @ip = @ip.succ
    puts "Standby Hostname: #{@hostname}-s"
    puts "Standby IP: #{@ip}",""
    build_vm( config, "#{@hostname}-s", "#{@ip}")
  end

  if (@segments > 0)
    puts "Segments Hosts: #{@segments}"
    # Segment Nodes:
    (1..@segments).each do |i|
      @ip = @ip.succ
      puts "[#{i}] Segment Hostname: #{@hostname}-#{i}"
      puts "[#{i}] Segment IP: #{@ip}",""
      build_vm( config, "#{@hostname}-#{i}", "#{@ip}")
    end
  end

  # Prepare the host that was provisioned
  # Provisioning
  config.vm.provision :shell, path: 'scripts/os.prep.sh', :args => ["#{@pivnet_token}", "#{@segments}"]
  # Developer Mode
  config.vm.provision :shell, path: 'scripts/go.build.sh', :args => ["#{@developer_mode}"]

end
