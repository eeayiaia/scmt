#!/bin/bash

# Input: CLUSTER_SUBNET

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

# Installing NFS
echo "Installing NFS"

apt-get install nfs-kernel-server --assume-yes
INSTALL_SUCCESS=$?

echo ""
write_line

if [[ $INSTALL_SUCCESS != 0 ]]; then
	echo "Failed to install NFS."
	exit 1
fi

# Filesystem that is to be exported needs to exist
mkdir /var/nfs

# Set ownership
chown nobody:nogroup /var/nfs

# Adding clients to the list that we will share with
echo "/var/nfs	$CLUSTER_SUBNET(rw,sync,no_subtree_check)" >> /etc/exports

# Create the nfs table
exportfs -a

# Start the service
service nfs-kernel-server start

