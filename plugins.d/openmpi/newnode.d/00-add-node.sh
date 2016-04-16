#!/bin/bash

# Inputs: NODE_IP, NODENAME

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1

check_invoked_by_scmt

echo "OpenMPI: adding node with hostname '$NODENAME' and IP '$NODE_IP'"

# Add new node to machine-file
OPENMPI_MACHINEFILE="/home/mpiuser/openmpi-machinefile"

if [[ ! -f "$OPENMPI_MACHINEFILE" ]]; then
	touch "$OPENMPI_MACHINEFILE"
	chown mpiuser:mpiuser "$OPENMPI_MACHINEFILE"
#else
#   Might be excessive to backup machinefile for each new node.
#	backup_file "$OPENMPI_MACHINEFILE"
fi

NUM_PROCS=$(nproc)
echo "$NODENAME	slots=$NUM_PROCS" >> "$OPENMPI_MACHINEFILE"

