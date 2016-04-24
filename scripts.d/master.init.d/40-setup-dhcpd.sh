#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

# Install DHCPD
apt-get install isc-dhcp-server --assume-yes

# Set up base DHCPD config
DHCPD_CONF=/etc/dhcp/dhcpd.conf

if [[ -f  "$DHCPD_CONF" ]]; then
	backup_file "$DHCPD_CONF"
fi

cp "$DIR/resources/baseDHCPD.conf" "$DHCPD_CONF"

