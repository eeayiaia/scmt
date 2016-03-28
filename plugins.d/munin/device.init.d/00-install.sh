#!/bin/sh

master_IP=$(awk '/^[[:space:]]*($|#)/{next} /master/{print $1; exit}' /etc/hosts)

if [[ $EUID -ne 0 ]]; then
	echo "This installer must be run with root rights." 1>&2
	exit 100
fi

apt-get install -y --force-yes munin-node

echo "cidr_allow $master_IP/32" >> /etc/munin/munin-node.conf

service munin-node restart