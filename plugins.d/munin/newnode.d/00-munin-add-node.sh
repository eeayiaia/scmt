#!/bin/bash
MASTER_INSTALLED=/usr/bin/munin-cron
IP_ADDRESS=$1
NODE_NAME=$2

if [ ! -x "$MASTER_INSTALLED" ]; then
  # if munin master is not installed (if /etc/munin is not a directory)
  echo install munin master first
  exit 1

fi

if [[ $EUID -ne 0 ]]; then
	echo "This installer must be run with root rights." 1>&2
	exit 100
fi

if [ -z $IP_ADDRESS ] || [ -z $NODE_NAME ]; then
	#if either ip address or node name does not exist
	echo please call this script "munin-add-node <node ip address> <node name>"
	exit 2
fi
#format what to add in the conf file
ADD_TO_CONF="[$NODE_NAME]
	address $IP_ADDRESS
	use_node_name yes"

#Backup config file
backup_file /etc/munin/munin.conf

#write to conf file
echo "$ADD_TO_CONF" >> /etc/munin/munin.conf

service apache2 restart
