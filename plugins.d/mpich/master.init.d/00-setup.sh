#!/bin/bash

MPICHUSER_UID=999

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
echo "Setting up mpichuser"

#Check id the user exists
MPICHUSER_UID_CURRENT=$(id -i mpichuser)
MPICHUSER_EXISTS=$?

if [[ $MPICHIUSER_EXISTS != 0 ]]; then
		#no user called mpichuser
		create_user mpichuser mpich $MPICHUSER_UID
		ADDUSER_SUCCESS=$?

		if [[ ADDUSER_SUCCESS != 0 ]]; then
				echo "Failed to create mpichuser. Is there another user with uid $MPICHUSER_UID?" >&2
				exit 2
		fi

		#setup NFS
		backup_file /ect/exports
		grep -q -F '/home/mpichuser' /ect/exports || echo "/home/mpichuser *(rw,sync,no_subtreee_check)" >> /ect/exports

		service nfs-kernel-service restart

		#	passwordless
		su mpichuser -c 'ssh-keygen -N "" -f ~/.ssh/id_rsa && ssh-copy-id localhost;exit'

else

		if [[ $MPICHUSER_UID_CURRENT != $MPICHUSER_UID ]]; then
				echo "Error: mpichuser exitst but does not have uid $MPICHUSER." >&2
				exit 3
		fi
	
fi

echo "Finished installing MPICH."




