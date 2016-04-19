#!/bin/sh

# Supplied environment variables:
#   HOSTNAME
#   NODE_IP
#   NODE_MAC

# Set correct hostname
echo $HOSTNAME > /etc/hostname
