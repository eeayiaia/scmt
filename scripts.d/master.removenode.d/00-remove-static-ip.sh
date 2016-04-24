#!/bin/bash

# Supplied environment variables:
#   NODENAME

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

DHCPD_CONF=/etc/dhcp/dhcpd.conf

echo "Removing $NODENAME from $DHCPD_CONF"

if [[ ! -f "$DHCPD_CONF" ]]; then
	echo "Error: $DHCPD_CONF does not exist." >&2
	exit 1
fi

backup_file "$DHCPD_CONF"

# Remove node from DHCPD conf
awk '
	BEGIN {flag=1}
	/^host '$NODENAME'/{flag=0;next}
	/\}/ && !flag{flag=1;next}
	flag
' "$DHCPD_CONF" > temp && mv -- temp "$DHCPD_CONF"

