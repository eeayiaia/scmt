#!/bin/bash

# Input: NODENAME

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1

check_invoked_by_scmt

if [[ ! $NODENAME ]]; then
	echo "Munin: removenode.d: Missing parameter 'NODENAME'. Exiting." >&2
	exit 2
fi

#Backup config file
backup_file /etc/munin/munin.conf

# Remove node from munin.conf
awk "
	BEGIN { flag=1 }
	/\[$NODENAME\]/ { flag=0; next }
	/^\[/ { flag=1 }
	flag
" /etc/munin/munin.conf > /etc/munin/munin.conf.tmp \
	&& mv /etc/munin/munin.conf.tmp /etc/munin/munin.conf

