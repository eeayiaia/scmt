#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1

check_invoked_by_scmt

echo "Installing OpenMPI"
write_line
echo ""

apt-get install openmpi-bin libopenmpi-dev --assume-yes
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
		echo "Failed to create mpiuser. Is there another user with uid " \
			"$MPIUSER_UID?" >&2
		exit 3
	fi

	# Set up NFS sharing of mpiuser's home directory
	backup_file /etc/exports
	grep -q -F '/home/mpiuser' /etc/exports \
		|| echo "/home/mpiuser *(rw,sync,no_subtree_check)" >> /etc/exports

	service nfs-kernel-server restart

	# Allow passwordless ssh between mpiusers
	if [[ ! -f /home/mpiuser/.ssh/id_rsa ]]; then
		su mpiuser -c \
			'ssh-keygen -N "" -f ~/.ssh/id_rsa && ssh-copy-id localhost;exit'
	fi
else
	if [[ $MPIUSER_UID_CURRENT != $MPIUSER_UID ]]; then
		echo "Error: mpiuser exists but does not have uid $MPIUSER_UID." >&2
		exit 4
	fi
fi

# Copy helper scripts
cp $DIR/../resources/run-with-openmpi.sh /home/mpiuser/
cp $DIR/../resources/compile-with-openmpi.sh /home/mpiuser/

chown mpiuser:mpiuser /home/mpiuser/run-with-openmpi.sh
chown mpiuser:mpiuser /home/mpiuser/compile-with-openmpi.sh

chmod +x /home/mpiuser/run-with-openmpi.sh
chmod +x /home/mpiuser/compile-with-openmpi.sh

# Copy test program
cp $DIR/../resources/mpi-hello-world.c /home/mpiuser/
chown mpiuser:mpiuser /home/mpiuser/mpi-hello-world.c

echo "Finished installing OpenMPI."

