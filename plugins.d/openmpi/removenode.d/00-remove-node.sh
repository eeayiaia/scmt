#!/bin/bash

# Inputs: NODE_IP, NODENAME

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh" || exit 1

check_root

echo "OpenMPI: removing node with hostname '$NODENAME' and IP '$NODE_IP'"

# Remove node from machine-file
OPENMPI_HOSTFILE="/home/mpiuser/openmpi-hostfile"

if [[ ! -f $OPENMPI_HOSTFILE ]]; then
	echo "Error: '$OPENMPI_HOSTFILE' does not exist." 1&>2
	exit 2
else
	backup_file $OPENMPI_HOSTFILE
fi

sed -i".bak" '/'$NODENAME'/d' $OPENMPI_HOSTFILE

