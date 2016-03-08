#!/bin/sh
if [[ $EUID -ne 0 ]]; then
	echo "This installer must be run with root rights." 1>&2
	exit 100
fi

apt-get install -y munin-node

#Get package including deps list url for debs
#apt-get --print-uris --yes install apache2 | grep ^\' | cut -d\' -f2
#not working if package already installed.