#!/bin/bash

# Get script directoy
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../script-utils/installer-utils.sh"
. "$DIR/openmpi-package-list.sh"

check_root

check_or_install_aptitude

# TODO: Check gcc/glibc version?

echo "Downloading OpenMPI packages..."
write_line

# Download necessary packages
aptitude download ${packages[@]}

write_line
echo "Finished downloading OpenMPI packages."

