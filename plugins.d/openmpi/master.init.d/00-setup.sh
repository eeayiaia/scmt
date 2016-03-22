#!/bin/bash

MPIUSER_UID=999

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh" || exit 1

check_root


echo "Installing OpenMPI"
write_line
echo ""

apt-get install openmpi-bin
INSTALL_SUCCESS=$?

echo ""
write_line

if [[ $INSTALL_SUCCESS != 0 ]]; then
	echo "Failed to install OpenMPI." >&2
	exit 2
fi

# Create mpiuser
echo "Setting up mpiuser"

# Check if mpiuser already exists
MPIUSER_UID_CURRENT=$(id -u mpiuser)
MPIUSER_EXISTS=$?

if [[ $MPIUSER_EXISTS != 0 ]]; then
	# No user called mpiuser
	create_user mpiuser mpi $MPIUSER_UID
	ADDUSER_SUCCESS=$?

	if [[ $ADDUSER_SUCCESS != 0 ]]; then
		echo "Failed to create mpiuser. Is there another user with uid $MPIUSER_UID?" >&2
		exit 3
	fi

	# Set up NFS sharing of mpiuser's home directory
	backup_file /etc/exports
	grep -q -F '/home/mpiuser' /etc/exports || echo "/home/mpiuser *(rw,sync,no_subtree_check)" >> /etc/exports

	service nfs-kernel-service restart

	# Allow passwordless ssh between mpiusers
	su mpiuser -c 'ssh-keygen -N "" -f ~/.ssh/id_rsa && ssh-copy-id localhost;exit'

else
	if [[ $MPIUSER_UID_CURRENT != $MPIUSER_UID ]]; then
		echo "Error: mpiuser exists but does not have uid $MPIUSER_UID." >&2
		exit 4
	fi	
fi

echo "Finished installing OpenMPI."

