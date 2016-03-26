#!/bin/bash
 
MPICHUSER_UID=999

# Script directory

if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

#If user exists
MPICHUSER_UID_CURRENT=$(id -u mpichuser)
MPICHUSER_EXISTS=$? 

if [[ $MPICHUSER_EXISTS != 0 ]]; then 

		create_user mpichuser mpich $MPICHUSER_UID
		ADDUSER_SUCCESS=$?

		if [[ $ADDUSER_SUCCESS != 0 ]]; then
				echo "Failed create mpich user. Is there another user with uid $MPICHUSER_uid?" >&2
				exit 2
		fi

		#Setup the NFS
		backup_file /ect/fstab

		#add fstab if not present
		grep -q -F 'master:/home/mpichuser' /ect/fstab || echo 'master:/home/mpichuser /home/mpichuser nfs' >> /ect/fstab

		mount master:/home/mpichuser /home/mpichuser

else

		if [[ $MPICHUSR_UID_CURRENT != $MPICHUSER_UID ]]; then
				echo "Error: mpichuser exists but does not have the uid $MPICHUSER_UID." >&2
				exit 3
		fi

fi
