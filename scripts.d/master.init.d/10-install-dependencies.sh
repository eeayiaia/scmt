#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

# Install aptitude
write_line
echo "Installing aptitude..."
apt-get install aptitude --assume-yes
write_line

# Install realpath
write_line
echo "Installing realpath..."
apt-get install realpath --assume-yes
write_line

