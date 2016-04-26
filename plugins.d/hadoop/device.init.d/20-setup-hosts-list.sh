#!/bin/bash

#Script runs as many times there are slave-nodes in the cluster 

$HOST_IP=$1
$NODE_NUMBER=$2


#Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1

#Set ip to /etc/hosts
echo "$HOST_IP hadoop-slave$2"
