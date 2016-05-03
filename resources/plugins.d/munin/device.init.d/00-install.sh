#!/bin/sh

# Include utils
. "/var/shared/utils.sh" || exit 1

if [[ $EUID -ne 0 ]]; then
	echo "This installer must be run with root rights." 1>&2
	exit 100
fi

apt-get install -y --force-yes munin-node


#Backup config file
backup_file /etc/munin/munin-node.conf

echo "cidr_allow $MASTER_IP/32" >> /etc/munin/munin-node.conf

service munin-node restart
