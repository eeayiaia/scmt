#!/bin/bash

# Input: MASTER_IP
# Set the hosts file to include the master node

. "$UTILS_PATH" || exit 1

check_invoked_by_scmt

assertIsSet MASTER_IP

backup_file /etc/hosts

egrep -q "^\s*\S+\s*master\s*$" /etc/hosts
UPDATE=$?

if [[ $UPDATE == 0 ]]; then
	echo "Updating master address in /etc/hosts"
	sed "s/.*\smaster/${MASTER_IP}    master/" -i /etc/hosts
else
	echo "Adding master to /etc/hosts"
	sed "$ a\
${MASTER_IP}    master" -i /etc/hosts
fi

