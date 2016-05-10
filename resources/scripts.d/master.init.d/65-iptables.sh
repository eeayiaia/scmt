#!/bin/bash

# Get script directory
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
. "$DIR/../utils.sh" || exit 1

#check_invoked_by_scmt

EXT=$NETWORK_INTERFACE_EXTERNAL
INT=$NETWORK_INTERFACE_INTERNAL

IPTABLES_CONF_FILE="/etc/iptables.rules"
IPTABLES_RULES="* nat\n
-A POSTROUTING -o $EXT -j MASQUERADE\n
COMMIT\n
\n
* accept all\n
-A INPUT -i lo -j ACCEPT\n
-A INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT\n
COMMIT"

echo "Setting up rules for iptables in file $IPTABLES_CONF_FILE"
echo -e $IPTABLES_RULES >| "$IPTABLES_CONF_FILE"
echo -e "\tenabled NAT on $EXT"
