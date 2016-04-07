#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

nodeIP=$1


#Backup config file
backup_file /etc/ganglia/gmetad.conf


python helpscript/regex.py $nodeIP