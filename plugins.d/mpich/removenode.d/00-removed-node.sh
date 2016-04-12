#!/bin/bash

# Input: NODE_IP, NODENAME

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

echo "MPICH: removing node with hostname '$NODENAME' and IP '$NODE_IP'"

MPICH_HOSTFILE="/home/mpichuser/mpich_hostfile"

if [[ ! -f $MPICH_HOSTFILE ]]; then
		echo "Error: '$MPICH_HOSTFILE' does not exist. " 1&>2
		exit 1
else
		backup_file $MPICH_HOSTFILE
fi

sed -i ".bak" '/'$NODENAME'/d' $MPICH_HOSTFILE

