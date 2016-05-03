#!/bin/bash

#Script runs once to remove hadoop-slave from /etc/hosts

$HOST_IP=$1
$NODE_ID=$2


#Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1

#Set ip to /etc/hosts
sed -e "s/$HOST_IP hadoop-slave$2/" 
