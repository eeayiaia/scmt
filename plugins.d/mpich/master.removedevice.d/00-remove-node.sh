#!/bin/bash

# Input: NODE_IP, NODENAME

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1

check_invoked_by_scmt

echo "MPICH: removing node with hostname '$NODENAME' and IP '$NODE_IP'"

MPICH_MACHINEFILE="/home/mpiuser/mpich-machinefile"

if [[ ! -f "$MPICH_MACHINEFILE" ]]; then
	echo "Error: '$MPICH_MACHINEFILE' does not exist. " 1&>2
	exit 1
#else
#	Might be excessive to backup machinefile for each removed node
#	backup_file $MPICH_MACHINEFILE
fi

sed -i".bak" '/'$NODENAME'/d' "$MPICH_MACHINEFILE"

