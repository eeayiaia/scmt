#!/bin/bash

MPIUSER_UID=999

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh" || exit 1

check_root

# Install OpenMPI
echo "Installing OpenMPI"
write_line
apt-get install openmpi-bin libopenmpi-dev --assume-yes
INSTALL_SUCCESS=$?

if [[ $INSTALL_SUCCESS != 0 ]]; then
	echo "Failed to install OpenMPI. Is approx configured correctly?" >&2
	exit 1
fi

# Check if mpiuser already exists
MPIUSER_UID_CURRENT=$(id -u mpiuser)
MPIUSER_EXISTS=$?

if [[ $MPIUSER_EXISTS != 0 ]]; then
	# No user called mpiuser
	create_user mpiuser mpi $MPIUSER_UID
	ADDUSER_SUCCESS=$?

	if [[ $ADDUSER_SUCCESS != 0 ]]; then
		echo "Failed to create mpiuser. Is there another user with uid $MPIUSER_UID?" >&2
		exit 2
	fi

	# Set up NFS mount for mpiuser home directory
	backup_file /etc/fstab

	# Add to fstab if not already present
	grep -q -F 'master:/home/mpiuser' /etc/fstab || echo 'master:/home/mpiuser /home/mpiuser nfs' >> /etc/fstab

	mount master:/home/mpiuser /home/mpiuser
else
	if [[ $MPIUSER_UID_CURRENT != $MPIUSER_UID ]]; then
		echo "Error: mpiuser exists but does not have uid $MPIUSER_UID." >&2
		exit 3
	fi	
fi

