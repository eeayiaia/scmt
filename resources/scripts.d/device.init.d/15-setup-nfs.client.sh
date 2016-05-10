#!/bin/bash

# Inputs: MASTER_IP

# Get script directory & include utils
. "$UTILS_PATH"

check_invoked_by_scmt

write_line
echo "Setting up NFS client"

# Install nfs
echo "Installing NFS..."
write_line
apt-get install rpcbind nfs-common --force-yes
INSTALL_SUCCESS=$?

write_line

if [[ $INSTALL_SUCCESS != 0 ]]; then
	echo "Failed to install NFS on node." >&2
	exit 1
fi

HOSTS_DENY=/etc/hosts.deny
echo "Adding rpcbind to $HOSTS_DENY"
backup_file "$HOSTS_DENY"

grep -q -F "rpcbind\s*:\s*ALL" "$HOSTS_DENY" \
	|| echo "rpcbind : ALL" >> "$HOSTS_DENY"

HOSTS_ALLOW=/etc/hosts.allow
echo "Adding allowed hosts to $HOSTS_ALLOW"
backup_file "$HOSTS_ALLOW"

grep -q -F "$MASTER_IP" "$HOSTS_ALLOW" \
	|| echo "$MASTER_IP" >> "$HOSTS_ALLOW"

# Create mounted shared folder
if [[ ! -d /var/shared ]]; then
	mkdir /var/shared
fi

# Mount the server shared directory to the local folder
echo "Mounting shared directory"
mount "$MASTER_IP:/var/nfs" /var/shared

# Adding mount to fstab
backup_file /etc/fstab

grep -q -F "$MASTER_IP:/var/nfs" /etc/fstab \
	|| echo "$MASTER_IP:/var/nfs /var/shared nfs" >> /etc/fstab

#Start services
echo "Starting services"
service portmap restart

echo "Finished setting up NFS client."

