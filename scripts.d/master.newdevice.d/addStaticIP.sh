#!/bin/sh
DHCPDPATH=/etc/dhcp/dhcpd.conf
# Supplied environment variables:
#   HOSTNAME
#   NODE_IP
#   NODE_MAC


echo -e "\n  host $HOSTNAME{" >> baseDHCPD.conf
echo "    hardware ethernet $NODE_MAC;" >> baseDHCPD.conf
echo "    fixed-address $NODE_IP;" >> baseDHCPD.conf
echo "  }" >> baseDHCPD.conf
cat baseDHCPD.conf > $DHCPDPATH
echo "}" >> $DHCPDPATH

