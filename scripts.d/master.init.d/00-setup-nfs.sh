#!/bin/bash

#If we need to check if NFS allready is installed uncomment following
#dpkg -l | grep nfs-kernel-server
#ALLREADY_INSTALLED=$?
#
#if [[ $ALLREADY_INSTALLED == 0]]; then
#		nfs is allready installed so exiting
#		exit 0
#fi


#Installing NFS
echo "Installing NFS"

sudo apt-get install nfs-kernel-server
INSTALL_SUCCESS=$?

echo ""
write_line

if [[ $INSTALL_SUCCESS != 0]]; then
		echo "Failed to install NFS."
		exit 1
fi

#Filesystem that is to be exported needs to exist
sudo mkdir /var/nfs

#set ownership
sudo chown nobody:nogroup /var/nfs

#Adding clients to the list that we will share with
echo "/var/nfs		10.46.0.101(rw)" >> /etc/exports
echo "/var/nfs		10.46.0.102(rw)" >> /etc/exports
echo "/var/nfs		10.46.0.103(rw)" >> /etc/exports

#Create the nfs table
sudo exportfs -a

#Start the service
sudo service nfs-kernel-server start










