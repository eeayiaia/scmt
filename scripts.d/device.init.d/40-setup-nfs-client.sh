#!/bin/bash

# Get script directory & include utils
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

#Install nfs
apt-get install portmap nfs-common --assume-yes
INSTALL_SUCCESS=$?

echo ""
write_line

if [[ $INSTALL_SUCCESS != 0 ]]; then
	echo "Failed to install NFS on node."
	exit 1
fi

#Adding rpcbind
echo "rpcbind : ALL" >> /etc/hosts.deny
 
#Adding allowed hosts
echo "10.46.0.1" >> /etc/hosts.allow

#Create mounted shared folder
if [[ ! -d /var/shared ]]; then
	mkdir /var/shared
fi

#Mount the server shared folder to the local folder
mount 10.46.0.1:/var/nfs /var/shared/

#Start services
/etc/init.d/portmap restart
/etc/init.d/nfs-common restart

