#!/bin/bash

# Supplied environment variables:
#   HOSTNAME

# Get script directory & include utils
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

# Set correct hostname
backup_file /etc/hostname
echo $NODE_HOSTNAME > /etc/hostname
echo "Set hostname to $NODE_HOSTNAME"
