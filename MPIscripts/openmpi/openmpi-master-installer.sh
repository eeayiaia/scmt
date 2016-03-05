#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../script-utils/installer-utils.sh"
. "$DIR/openmpi-package-list.sh"

check_root

for i in "${packages[@]}"
do
	install_pkg $i
done

echo "Finished installing OpenMPI."

