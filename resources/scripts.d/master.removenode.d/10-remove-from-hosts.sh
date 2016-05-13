#!/bin/bash

# Remove node from hosts file
# Supplied environment variables:
#  NODENAME

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

HOSTSFILE=/etc/hosts

backup_file $HOSTSFILE

echo "Removing $NODENAME from $HOSTSFILE..."

# Remove if present
sed '/\s'$NODENAME'/d' $HOSTSFILE > temp && mv temp /etc/hosts

