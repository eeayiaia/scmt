#!/bin/bash

#Script runs once to add hadoop-master to /etc/hosts

$HOST_IP=$1



#Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1

#Set ip to /etc/hosts
echo "$HOST_IP hadoop-master"
