#!/bin/bash

# Supplied environment variables:
#   NODE_HOSTNAME

. "$UTILS_PATH"

check_invoked_by_scmt

backup_file "/etc/hostname"

assertIsSet NODE_HOSTNAME

# Find old hostname
OLDNAME=$(cat /etc/hostname)

# Set correct hostname
echo $NODE_HOSTNAME >| /etc/hostname
echo "Set hostname to $NODE_HOSTNAME"

# Update hostfile to include new hostname
if [[ $OLDNAME == "" ]]; then
	# This should of course not happen, but just to be sure...
	echo "127.0.0.1	${NODE_HOSTNAME}" >> /etc/hosts
else
	echo "Updating hostname in /etc/hosts"
	sed "s/.*${OLDNAME}/127.0.0.1    ${NODE_HOSTNAME}/" -i /etc/hosts
fi

