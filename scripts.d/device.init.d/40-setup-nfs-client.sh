#!/bin/bash


#Install nfs
sudo apt-get install portmap nfs-common
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
sudo mkdir /var/shared

#Mount the server shared folder to the local folder
sudo mount 10.46.0.1:/var/nfs /var/shared/


#Start services
sudo /etc/init.d/portmap restart
sudo /etc/init.d/nfs-common restart



