#!/bin/bash

MPIUSER_UID=999

# Script directory

if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

echo "Installing MPICH"
write_line
echo""

apt-get install mpich2
INSTALL_SUCCESS=$?

echo""
write_line

if [[ $INSTALL_SUCCESS != 0 ]]; then
		echo "Failed to install MPICH." >&2
		exit 1
fi

#Create user
echo "Setting up mpiuser"

#Check id the user exists
MPIUSER_UID_CURRENT=$(id -i mpichuser)
MPIUSER_EXISTS=$?

if [[ $MPIUSER_EXISTS != 0 ]]; then
		#no user called mpichuser
		create_user mpichuser mpich $MPIUSER_UID
		ADDUSER_SUCCESS=$?

		if [[ ADDUSER_SUCCESS != 0 ]]; then
				echo "Failed to create mpichuser. Is there another user with uid $MPIUSER_UID?" >&2
				exit 2
		fi

		#setup NFS
		backup_file /ect/exports
		grep -q -F '/home/mpiuser' /ect/exports || echo "/home/mpiuser *(rw,sync,no_subtreee_check)" >> /ect/exports

		service nfs-kernel-service restart

		#	passwordless
		su mpichuser -c 'ssh-keygen -N "" -f ~/.ssh/id_rsa && ssh-copy-id localhost;exit'

else

		if [[ $MPIUSER_UID_CURRENT != $MPIUSER_UID ]]; then
				echo "Error: mpiuser exitst but does not have uid $MPIUSER." >&2
				exit 3
		fi
	
fi

echo "Finished installing MPICH."




