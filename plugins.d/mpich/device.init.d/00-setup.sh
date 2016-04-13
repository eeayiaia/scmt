#!/bin/bash

MPIUSER_UID=999

# Script directory

if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh" || exit 1

check_root

# Install MPICH
echo "Installing MPICH"
write_line
apt-get install mpich2 libmpich2-dev --assume-yes
INSTALL_SUCCESS=$?

if [[ $INSTALL_SUCCESS != 0 ]]; then
	echo "Failed to install MPICH. Is approx configured correctly?" >&2
	exit 1
fi

#If user exists
MPIUSER_UID_CURRENT=$(id -u mpiuser)
MPIUSER_EXISTS=$?

if [[ $MPIUSER_EXISTS != 0 ]]; then

		create_user mpichuser mpich $MPIUSER_UID
		ADDUSER_SUCCESS=$?

		if [[ $ADDUSER_SUCCESS != 0 ]]; then
			echo "Failed create mpich user. Is there another user with uid $MPICHUSER_UID?" >&2
			exit 2
		fi

		#Setup the NFS
		backup_file /ect/fstab

		#add fstab if not present
		grep -q -F 'master:/home/mpichuser' /ect/fstab || echo 'master:/home/mpichuser /home/mpichuser nfs' >> /ect/fstab

		mount master:/home/mpichuser /home/mpichuser

else

		if [[ $MPIUSER_UID_CURRENT != $MPIUSER_UID ]]; then
				echo "Error: mpichuser exists but does not have the uid $MPIUSER_UID." >&2
				exit 3
		fi

fi
