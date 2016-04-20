#!/bin/sh

# Input: MASTER_IP
# Set the hosts file to include the master node

egrep -q "^master\s" /etc/hosts \
	&& sed "s/master/$MASTER_IP    master/" -i /etc/hosts \
	|| sed "$ a\$MASTER_IP    master" -i /etc/hosts

