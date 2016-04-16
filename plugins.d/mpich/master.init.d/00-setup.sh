#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1

check_invoked_by_scmt

echo "Installing MPICH"
write_line
echo ""

apt-get install mpich2 libmpich2-dev --assume-yes
INSTALL_SUCCESS=$?

echo ""
write_line

if [[ $INSTALL_SUCCESS != 0 ]]; then
	echo "Failed to install MPICH." >&2
	exit 1
fi

# Create mpiuser
echo "Setting up mpiuser"

# Check if mpiuser exists
MPIUSER_UID_CURRENT=$(id -u mpiuser)
MPIUSER_EXISTS=$?

if [[ $MPIUSER_EXISTS != 0 ]]; then
	# No user called mpiuser
	create_user mpiuser mpi $MPIUSER_UID
	ADDUSER_SUCCESS=$?

	if [[ $ADDUSER_SUCCESS != 0 ]]; then
		echo "Failed to create mpiuser. Is there another user with uid " \
			"$MPIUSER_UID?" >&2
		exit 2
	fi

	# Set up NFS sharing of mpiuser's home directory
	backup_file /etc/exports
	grep -q -F '/home/mpiuser' /etc/exports \
		|| echo "/home/mpiuser *(rw,sync,no_subtree_check)" >> /etc/exports

	service nfs-kernel-service restart

	# Allow passwordless ssh between mpiusers
	if [[ ! -f /home/mpiuser/.ssh/id_rsa ]]; then
		su mpiuser -c \
			'ssh-keygen -N "" -f ~/.ssh/id_rsa && ssh-copy-id localhost;exit'
	fi
else
	if [[ $MPIUSER_UID_CURRENT != $MPIUSER_UID ]]; then
		echo "Error: mpiuser exists but does not have uid $MPIUSER_UID." >&2
		exit 3
	fi
fi

# Copy helper scripts
cp $DIR/../resources/run-with-mpich.sh /home/mpiuser/
cp $DIR/../resources/compile-with-mpich.sh /home/mpiuser/

chown mpiuser:mpiuser /home/mpiuser/run-with-mpich.sh
chown mpiuser:mpiuser /home/mpiuser/compile-with-mpich.sh

chmod +x /home/mpiuser/run-with-mpich.sh
chmod +x /home/mpiuser/compile-with-mpich.sh

# Copy example MPI program
cp $DIR/../resources/mpi-hello-world.c /home/mpiuser/
chown mpiuser:mpiuser /home/mpiuser/mpi-hello-world.c

echo "Finished installing MPICH."

