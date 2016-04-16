#!/bin/bash
MASTER_INSTALLED=/usr/bin/munin-cron

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1

check_invoked_by_scmt

if [ ! -x "$MASTER_INSTALLED" ]; then
  # if munin master is not installed (if /etc/munin is not a directory)
  echo install munin master first
  exit 1

fi

if [ -z $NODE_IP ] || [ -z $NODENAME ]; then
	#if either ip address or node name does not exist
	echo please call this script "munin-add-node <node ip address> <node name>"
	exit 2
fi
#format what to add in the conf file
ADD_TO_CONF="[$NODENAME]
	address $NODE_IP
	use_node_name yes"

#Backup config file
backup_file /etc/munin/munin.conf

#write to conf file
echo "$ADD_TO_CONF" >> /etc/munin/munin.conf

service apache2 restart
