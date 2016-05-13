#!/bin/bash

# Add node to hosts file
# Supplied environment variables:
#  NODENAME
#  NODE_IP

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

HOSTSFILE=/etc/hosts

backup_file $HOSTSFILE

# Add new node to hosts file
echo "Adding ($NODE_IP, $NODENAME) to $HOSTSFILE..."

# Remove if already present, then append new line
sed '/\s'$NODENAME'/d' $HOSTSFILE > temp && mv temp /etc/hosts

echo "$NODE_IP	$NODENAME" >> $HOSTSFILE

