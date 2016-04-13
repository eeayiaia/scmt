#!/bin/bash

# Input: NODE_IP, NODENAME

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

echo "MPICH: adding node with hostname '$NODENAME' and IP '$NODE_IP'"

MPICH_MACHINEFILE="/home/mpiuser/mpich-machinefile"

if [[ ! -f $MPICH_MACHINEFILE ]]; then
		touch $MPICH_MACHINEFILE
		chown mpiuser:mpiuser $MPICH_MACHINEFILE
else
		backup_file $MPICH_MACHINEFILE
fi

echo "$NODENAME:4" >> $MPICH_MACHINEFILE

