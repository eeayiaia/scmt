#!/bin/bash

NODE_IP=$1
NODE_NAME=$2


DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

echo "MPICH: adding node with hostname '$NODE_NAME' and IP 'NODE_IP'"

MPICH_HOSTFILE="/home/mpich-hostfile"

if [[ ! -f $MPICH_HOSTFILE ]]; then
		touch $MPICH_HOSTFILE
		chown mpich:mpichuser $MPICHUSER_HOSTFILE
else
		backup_file $MPICH_HOSTFILE
fi

echo "$NODE_NAME		slots=4" >> $MPICH_HOSTFILE

