#!/bin/bash

# Input: MASTER_IP
# Set the hosts file to include the master node

. "$UTILS_PATH"

check_invoked_by_scmt

backup_file /etc/hosts

egrep -q "^master\s" /etc/hosts \
	&& sed "s/master/${MASTER_IP}    master/" -i /etc/hosts \
	|| sed "$ a\
${MASTER_IP}    master" -i /etc/hosts

