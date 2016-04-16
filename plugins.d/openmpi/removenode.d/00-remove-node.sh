#!/bin/bash

# Inputs: NODE_IP, NODENAME

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1

check_invoked_by_scmt

echo "OpenMPI: removing node with hostname '$NODENAME' and IP '$NODE_IP'"

# Remove node from machine-file
OPENMPI_MACHINEFILE="/home/mpiuser/openmpi-machinefile"

if [[ ! -f "$OPENMPI_MACHINEFILE" ]]; then
	echo "Error: '$OPENMPI_MACHINEFILE' does not exist." 1&>2
	exit 2
#else
#	Might be excessive to backup machinefile for each removed node
#	backup_file "$OPENMPI_MACHINEFILE"
fi

sed -i".bak" '/'$NODENAME'/d' "$OPENMPI_MACHINEFILE"

