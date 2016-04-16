#!/bin/bash

# Input: NODE_IP

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1

check_invoked_by_scmt

if [ -z $NODE_IP ]; then
	#if either ip address or node name does not exist
	echo please call this script "munin-remove-node <node ip address>"
	exit 2
fi

#Backup config file
backup_file /etc/munin/munin.conf

python helpscript/removenode.py "$NODE_IP"
