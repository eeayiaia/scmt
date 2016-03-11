#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

IP_ADDRESS=$1

check_root

if [ -z $IP_ADDRESS ]; then
	#if either ip address or node name does not exist
	echo please call this script "munin-remove-node <node ip address>"
	exit 2
fi

python helpscript/removenode.py "$IP_ADDRESS"

#format what to add in the conf file
ADD_TO_CONF="[$NODE_NAME]
	address $IP_ADDRESS
	use_node_name yes"

#write to conf file
echo "$ADD_TO_CONF" >> /etc/munin/munin.conf
