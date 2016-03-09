#!/bin/bash
NODE_IP=$1
NODE_NAME=$2

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

echo "OpenMPI: removing node with hostname '$NODE_NAME' and IP '$NODE_IP'"

# Remove node from machine-file
OPENMPI_HOSTFILE="/home/mpiuser/openmpi-hostfile"

if [[ ! -f $OPENMPI_HOSTFILE ]]; then
	echo "Error: '$OPENMPI_HOSTFILE' does not exist." 1&>2
	exit 1
else
	backup_file $OPENMPI_HOSTFILE
fi

sed -i".bak" '/'$NODE_NAME'/d' $OPENMPI_HOSTFILE

