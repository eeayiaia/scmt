#!/bin/bash

# Get script directory
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

EXT=$NETWORK_INTERFACE_EXTERNAL
INT=$NETWORK_INTERFACE_INTERNAL
INTERFACES="/etc/network/interfaces"

backup_file "$INTERFACES"

CONF="\n
auto lo\n
iface lo inet loopback\n
\n
auto $EXT\n
iface $EXT inet dhcp\n
\n
auto $INT\n
iface $INT inet static\n
\taddress $MASTER_IP\n
\tbroadcast $CLUSTER_BROADCAST_IP\n
\tnetmask $CLUSTER_SUBNET_MASK\n
\tpost-up /sbin/iptables-restore < /etc/iptables.rules\n"

function setup_interfaces {
  echo -e $CONF >| "$INTERFACES"
  echo "Wrote interfaces configuration"
  echo -e "\t$EXT as external interface"
  echo -e "\t$INT as internal interface"
}

function setup_ip_forward {
  egrep "^net\.ipv4\.ip_forward\=1$" /etc/sysctl.conf 
  FOUND=$?

  if [[ $FOUND == 0 ]]; then
    echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
    echo "Set net.ipv4.ip_forward=1"
  fi
}

setup_interfaces
setup_ip_forward
