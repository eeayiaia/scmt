#!/bin/bash

# Supplied environment variables:
#   NODENAME
#   NODE_IP
#   NODE_MAC

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

DHCPD_CONF=/etc/dhcp/dhcpd.conf

echo "Adding $NODENAME with IP $NODE_IP and MAC $NODE_MAC to $DHCPD_CONF"

if [[ ! -f "$DHCPD_CONF" ]]; then
	echo "Error: $DHCPD_CONF does not exist." >&2
	exit 1
fi

backup_file "$DHCPD_CONF"

# If nodename for some reason already exists in dhcpd conf, remove it
awk '
	BEGIN {flag=1}
	/^host '$NODENAME'/{flag=0;next}
	/\}/ && !flag{flag=1;next}
	flag
' "$DHCPD_CONF" > temp && mv -- temp "$DHCPD_CONF"

echo "
host $NODENAME {
  hardware ethernet $NODE_MAC;
  fixed-address $NODE_IP;
}" >> "$DHCPD_CONF"

