#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

# Install approx
APPROX_PATH=$(which approx)

if [[ ! "$APPROX_PATH" ]]; then
	echo "approx not found, installing..."
	write_line
	apt-get install -y approx
	INSTALL_SUCCESS=$?
	write_line

	if [[ $INSTALL_SUCCESS != 0 ]]; then
		echo "Failed to install approx. Is the master node connected to the "\
			"internet?"
		exit 2
	fi
else
	echo "approx found, skipping installation"
fi

# Apply approx config
APPROX_CONF_SOURCE="$DIR/../../config/approx/approx.conf"
APPROX_CONF_TARGET="/etc/approx/approx.conf"

echo "Applying approx configuration..."

if [[ ! -f "$APPROX_CONF_SOURCE" ]]; then
	echo "File missing: '<scmt-root>/config/approx/approx.conf'. Failed to set"\
		" up approx."
	exit 3
fi

# First backup old config if it exists
if [[ -f "$APPROX_CONF_TARGET" ]]; then backup_file "$APPROX_CONF_TARGET"; fi

cp -rf "$APPROX_CONF_SOURCE" "$APPROX_CONF_TARGET"

# Restart initd to make approx config take effect
/etc/init.d/openbsd-inetd restart

echo "Finished setting up approx."

write_line

