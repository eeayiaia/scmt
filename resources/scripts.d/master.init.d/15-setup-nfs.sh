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

# Adding cluster subnet to the exports list
echo "Configuring /etc/exports"

EXPORTS=/etc/exports
backup_file "$EXPORTS"

egrep -q "/var/nfs\s+$CLUSTER_SUBNET" "$EXPORTS" \
	|| echo "/var/nfs	$CLUSTER_SUBNET(rw,sync,no_subtree_check)" >> "$EXPORTS"

# Link the correct directories into nfs
SCRIPTS_PATH=$(realpath "$DIR/../../scripts.d")
PLUGINS_PATH=$(realpath "$DIR/../../plugins.d")
CONFIGS_PATH=$(realpath "$DIR/../../configs")
UTILS_PATH=$(realpath "$DIR/../utils.sh")

SCRIPTS_TARGET="/var/nfs/scripts.d"
PLUGINS_TARGET="/var/nfs/plugins.d"
CONFIGS_TARGET="/var/nfs/configs"
UTILS_TARGET="/var/nfs/utils.sh"

[[ -f "$SCRIPTS_TARGET" ]] && delete_file "$SCRIPTS_TARGET"
[[ -f "$PLUGINS_TARGET" ]] && delete_file "$PLUGINS_TARGET"
[[ -f "$CONFIGS_TARGET" ]] && delete_file "$CONFIGS_TARGET"
[[ -f "$UTILS_TARGET" ]] && delete_file "$UTILS_TARGET"

cp "$SCRIPTS_PATH" "$SCRIPTS_TARGET"
cp "$PLUGINS_PATH" "$PLUGINS_TARGET"
cp "$CONFIGS_PATH" "$CONFIGS_TARGET"
cp "$UTILS_PATH" "$UTILS_TARGET"

# Create the nfs table
echo "Running 'exportfs -a'"
exportfs -a

# Start the service
echo "Starting nfs-kernel-server"
service nfs-kernel-server start

