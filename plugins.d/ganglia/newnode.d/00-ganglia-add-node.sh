#!/bin/bash

# Inputs: NODE_IP

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

#Backup config file
backup_file /etc/ganglia/gmetad.conf

sed -i "/^data_source/ s/$/ $NODE_IP/" /etc/ganglia/gmetad.conf
