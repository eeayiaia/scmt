#!/bin/bash

#Script runs as many times there are slave-nodes in the cluster 



#Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1

#Set ip to /etc/hosts
echo "$HOST_IP hadoop-slave$HOST_IP"
