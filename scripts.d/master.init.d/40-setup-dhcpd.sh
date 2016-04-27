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

echo "
subnet $CLUSTER_SUBNET_IP netmask $CLUSTER_SUBNET_MASK {
  range $DEVICE_IP_RANGE_BEGIN $DEVICE_IP_RANGE_END;
  option routers $MASTER_IP;
  option broadcast-address $CLUSTER_BROADCAST_IP;

  default-lease-time $DHCPD_LEASE_TIME_DEFAULT;
  max-lease-time $DHCPD_LEASE_TIME_MAX;
}

# Devices in cluster
" >> "$DHCPD_CONF"

