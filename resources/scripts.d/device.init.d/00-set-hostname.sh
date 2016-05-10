#!/bin/bash

# Supplied environment variables:
#   NODE_HOSTNAME

. "$UTILS_PATH"

check_invoked_by_scmt

backup_file "/etc/hostname"

# Set correct hostname
echo $NODE_HOSTNAME >| /etc/hostname
echo "Set hostname to $NODE_HOSTNAME"

