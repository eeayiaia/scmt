#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

function write_line(){
	echo "--------------------------------------------------------------------------------"
}

# Check if root
if [[ $EUID != 0 ]]; then
	echo "setup_approx.sh must be run with root privileges." 1>&2
	exit 1
fi

# Install approx
approx_path=$(which approx)

if [[ ! $approx_path ]]; then
	echo "approx not found, installing..."
	write_line
	apt-get install approx
	install_success=$?
	write_line

	if [[ $install_success != 0 ]]; then
		echo "Failed to install approx."
		exit 1
	fi
else
	echo "approx found."
fi

# Apply approx config
cp -rf /etc/approx/approx.conf /etc/approx/approx.conf.backup
cp -rf $DIR../../config/approx/approx.conf /etc/approx.conf

# Restart initd to make approx config take effect
/etc/init.d/openbsd-initd restart

