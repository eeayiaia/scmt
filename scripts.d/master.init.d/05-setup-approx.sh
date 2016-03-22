#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../plugins.d/.script-utils/installer-utils.sh" || exit 1

check_root

# Install approx
approx_path=$(which approx)

if [[ ! $approx_path ]]; then
	echo "approx not found, installing..."
	write_line
	apt-get install -y approx
	install_success=$?
	write_line

	if [[ $install_success != 0 ]]; then
		echo "Failed to install approx. Is the master node connected to the internet?"
		exit 2
	fi
else
	echo "approx found, skipping installation"
fi

# Apply approx config
APPROX_CONF_SOURCE="$DIR/../../config/approx/approx.conf"
APPROX_CONF_TARGET="/etc/approx/approx.conf"

echo "Applying approx configuration..."

if [[ ! -f $APPROX_CONF_SOURCE ]]; then
	echo "File missing: '<scmt-root>/config/approx/approx.conf'. Failed to set up approx."
	exit 3
fi

if [[ -f $APPROX_CONF_TARGET ]]; then backup_file $APPROX_CONF_TARGET; fi

cp -rf $APPROX_CONF_SOURCE $APPROX_CONF_TARGET

# Restart initd to make approx config take effect
/etc/init.d/openbsd-inetd restart

echo "Finished setting up approx."

write_line

