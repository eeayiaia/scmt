#!/bin/bash

# Include utils
. "/var/shared/utils.sh" || exit 1

# Include config
. "/var/shared/plugins.d/mpich/resources/config" || exit 1

check_invoked_by_scmt

# Install MPICH
echo "Installing MPICH"
write_line
apt-get install mpich2 libmpich2-dev --assume-yes
INSTALL_SUCCESS=$?

if [[ $INSTALL_SUCCESS != 0 ]]; then
	echo "Failed to install MPICH. Is approx configured correctly?" >&2
	exit 1
fi

# Check if mpiuser exists
MPIUSER_UID_CURRENT=$(id -u mpiuser)
MPIUSER_EXISTS=$?

if [[ $MPIUSER_EXISTS != 0 ]]; then
	# Create mpiuser
	create_user mpiuser mpi $MPIUSER_UID
	ADDUSER_SUCCESS=$?

	if [[ $ADDUSER_SUCCESS != 0 ]]; then
		echo "Failed create mpich user. Is there another user with uid " \
			"$MPIUSER_UID?" >&2
		exit 2
	fi

	# Setup NFS mount for mpiuser home directory
	backup_file /etc/fstab

	# Add to fstab if not already present
	grep -q -F 'master:/home/mpiuser' /etc/fstab \
		|| echo 'master:/home/mpiuser /home/mpiuser nfs' >> /etc/fstab

	mount master:/home/mpiuser /home/mpiuser
else
	if [[ $MPIUSER_UID_CURRENT != $MPIUSER_UID ]]; then
		echo "Error: mpiuser exists but does not have uid '$MPIUSER_UID'." >&2
		exit 3
	fi
fi

