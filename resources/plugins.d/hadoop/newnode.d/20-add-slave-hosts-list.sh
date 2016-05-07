#!/bin/bash

#Script needs to add the slave in /etc/hosts on master node

#Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1

#Set ip to /etc/hosts in the master node
echo "$NODE_IP hadoop-slave"
