#!/bin/bash

# Inputs: NODE_IP, NODENAME

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh" || exit 1

check_root

echo "OpenMPI: adding node with hostname '$NODENAME' and IP '$NODE_IP'"

# Add new node to machine-file
OPENMPI_HOSTFILE="/home/mpiuser/openmpi-hostfile"

if [[ ! -f $OPENMPI_HOSTFILE ]]; then
	touch $OPENMPI_HOSTFILE
	chown mpiuser:mpiuser $OPENMPI_HOSTFILE
else
	backup_file $OPENMPI_HOSTFILE
fi

# TODO: number of slots should not be hardcoded
NUM_PROCS=$(nproc)
echo "$NODENAME	slots=$NUM_PROCS" >> $OPENMPI_HOSTFILE

