#!/bin/sh

# Set the hosts file to include the masternode

echo 'master    10.46.0.1' | sudo tee -a /etc/hosts
