#!/bin/bash

# Set the nameserver .. (enforces)

. "$UTILS_PATH" || exit 1

check_invoked_by_scmt

echo "8.8.8.8" >| /etc/resolv.conf
