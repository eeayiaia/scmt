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
  
  on commit {
    set clip = binary-to-ascii(10, 8, \".\", leased-address);
    set clhw = concat (
      suffix (concat (\"0\", binary-to-ascii (16, 8, \"\", substring(hardware,1,1))),2), \":\",
      suffix (concat (\"0\", binary-to-ascii (16, 8, \"\", substring(hardware,2,1))),2), \":\",
      suffix (concat (\"0\", binary-to-ascii (16, 8, \"\", substring(hardware,3,1))),2), \":\",
      suffix (concat (\"0\", binary-to-ascii (16, 8, \"\", substring(hardware,4,1))),2), \":\",
      suffix (concat (\"0\", binary-to-ascii (16, 8, \"\", substring(hardware,5,1))),2), \":\",
      suffix (concat (\"0\", binary-to-ascii (16, 8, \"\", substring(hardware,6,1))),2));

    execute(\"/usr/bin/scmt\", \"register-device\",  clhw, clip);
  }

  default-lease-time $DHCPD_LEASE_TIME_DEFAULT;
  max-lease-time $DHCPD_LEASE_TIME_MAX;
}

# Devices in cluster
" >> "$DHCPD_CONF"

