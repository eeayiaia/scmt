#!/bin/sh

# Set the hosts file to include the masternode

MASTER_IP="10.46.0.1"

egrep -q "^master\s" /etc/hosts \
	&& sed "s/^master.*/master    $MASTER_IP/" -i /etc/hosts \
	|| sed "$ a\master    $MASTER_IP" -i /etc/hosts

