# -*- mode: ruby -*-
# vi: set ft=ruby :


## In Skallago, Vagrant is used to spin up a docker host.


VAGRANTFILE_API_VERSION = "2"

def MapPort(vm, port)
    vm.network "forwarded_port", guest: port, host_ip: "127.0.0.1", host: port
end

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

	config.vm.define :docko do |docko|
	    docko.vm.define "docko"
	    docko.vm.box = "yungsang/boot2docker"
	    docko.vm.box_version = "= 1.3.1"

		MapPort(docko.vm, 8000) # api service

		docko.vm.network "private_network", ip: "192.168.33.10"
		docko.vm.synced_folder "~", ENV['HOME'], type: "nfs"

		docko.vm.provider "virtualbox" do |v|
			v.memory = 4096
			v.cpus = 4
		end

		# The following two provisions were taken from the yungsang/boot2docker
		# documentation here: https://vagrantcloud.com/yungsang/boxes/boot2docker

		# Fix busybox/udhcpc issue
		docko.vm.provision :shell do |s|
			s.inline = <<-EOT
				if ! grep -qs ^nameserver /etc/resolv.conf; then
					sudo /sbin/udhcpc
				fi
				cat /etc/resolv.conf
			EOT
		end

		# Adjust datetime after suspend and resume
		docko.vm.provision :shell do |s|
			s.inline = <<-EOT
				sudo /usr/local/bin/ntpclient -s -h pool.ntp.org
				date
			EOT
		end
	end

end
