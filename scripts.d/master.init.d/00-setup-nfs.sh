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

apt-get install nfs-kernel-server
INSTALL_SUCCESS=$?

echo ""
write_line

if [[ $INSTALL_SUCCESS != 0]]; then
		echo "Failed to install NFS."
		exit 1
fi

#Filesystem that is to be exported needs to exist
#mkdir...  filesys goes here

#Need to mount filesys with permission 777
#mount --bind /local_filesys /export_filesys

#Need to write filesys into fstab else needs to be written every boot
echo "/local_filesys		/remote_filesys		none		bind		0		0" >> /ect/fstab





