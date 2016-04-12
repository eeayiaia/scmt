#!/bin/bash

# Input: NODE_IP, NODENAME

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

echo "MPICH: adding node with hostname '$NODENAME' and IP 'NODE_IP'"

MPICH_HOSTFILE="/home/mpich-hostfile"

if [[ ! -f $MPICH_HOSTFILE ]]; then
		touch $MPICH_HOSTFILE
		chown mpich:mpichuser $MPICHUSER_HOSTFILE
else
		backup_file $MPICH_HOSTFILE
fi

echo "$NODENAME:4" >> $MPICH_HOSTFILE

