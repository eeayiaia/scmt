#!/bin/bash

# Input: CLUSTER_SUBNET

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

# Installing NFS
echo "Installing NFS"

apt-get install nfs-kernel-server --assume-yes
INSTALL_SUCCESS=$?

echo ""
write_line

if [[ $INSTALL_SUCCESS != 0 ]]; then
	echo "Failed to install NFS."
	exit 1
fi

# Filesystem that is to be exported needs to exist
[[ ! -d /var/nfs ]] && mkdir /var/nfs

# Set ownership
chown nobody:nogroup /var/nfs

# Adding clients to the list that we will share with
EXPORTS=/etc/exports
backup_file "$EXPORTS"
echo "/var/nfs	$CLUSTER_SUBNET(rw,sync,no_subtree_check)" >> "$EXPORTS"

# Link the correct directories into nfs
SCRIPTS_PATH=$(realpath '$DIR/../../scripts.d')
PLUGINS_PATH=$(realpath '$DIR/../../plugins.d')
CONFIGS_PATH=$(realpath '$DIR/../../configs')
UTILS_PATH=$(realpath '$DIR/../utils.sh')

[[ -f "$SCRIPTS_PATH" ]] && delete_file "$SCRIPTS_PATH"
[[ -f "$PLUGINS_PATH" ]] && delete_file "$PLUGINS_PATH"
[[ -f "$CONFIGS_PATH" ]] && delete_file "$CONFIGS_PATH"
[[ -f "$UTILS_PATH" ]] && delete_file "$UTILS_PATH"

ln -sf "$SCRIPTS_PATH" "/var/nfs/scripts.d"
ln -sf "$PLUGINS_PATH" "/var/nfs/plugins.d"
ln -sf "$CONFIGS_PATH" "/var/nfs/configs"
ln -sf "$UTILS_PATH" "/var/nfs/utils.sh"

# Create the nfs table
exportfs -a

# Start the service
service nfs-kernel-server start

