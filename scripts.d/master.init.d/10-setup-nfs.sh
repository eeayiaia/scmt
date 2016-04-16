#!/bin/bash

#If we need to check if NFS allready is installed uncomment following
#dpkg -l | grep nfs-kernel-server
#ALLREADY_INSTALLED=$?
#
#if [[ $ALLREADY_INSTALLED == 0]]; then
#		nfs is allready installed so exiting
#		exit 0
#fi

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

#Installing NFS
echo "Installing NFS"

apt-get install nfs-kernel-server --assume-yes
INSTALL_SUCCESS=$?

echo ""
write_line

if [[ $INSTALL_SUCCESS != 0 ]]; then
		echo "Failed to install NFS."
		exit 1
fi

#Filesystem that is to be exported needs to exist
mkdir /var/nfs

#set ownership
chown nobody:nogroup /var/nfs

#Adding clients to the list that we will share with
# TODO: Make sure subnet is correct
echo "/var/nfs	10.46.0.0/24(rw,sync,no_subtree_check)" >> /etc/exports

#Create the nfs table
exportfs -a

#Start the service
service nfs-kernel-server start

