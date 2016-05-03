#!/bin/bash

# Get script directory
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

# Set up base DHCPD config
APPARMOR_CONF="/etc/apparmor.d/local/usr.sbin.dhcpd"
SCMT_PATH=$(readlink -f $DIR/../../..)

if [ -d "/etc/apparmor.d" ]
then
	echo "Apparmor exists."

	if [[ -f  "$DHCPD_CONF" ]]; then
		backup_file "$DHCPD_CONF"
	fi

	echo "Adding scmt binary to apparmor config."
    grep -q "$SCMT_PATH/scmt Uxr," "$APPARMOR_CONF" || echo "$SCMT_PATH/scmt Uxr," >> "$APPARMOR_CONF"
    grep -q "$SCMT_PATH/run-scmt.sh Uxr," "$APPARMOR_CONF" || echo "$SCMT_PATH/run-scmt.sh Uxr," >> "$APPARMOR_CONF"
    grep -q "/bin/dash Ux," "$APPARMOR_CONF" || echo "/bin/dash Ux," >> "$APPARMOR_CONF"

	echo "reloading config."
	apparmor_parser -r /etc/apparmor.d/usr.sbin.dhcpd
else
    echo "Apparmor doesn't exist."
fi
