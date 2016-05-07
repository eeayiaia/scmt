#!/bin/bash

#Script runs once to add hadoop-master to /etc/hosts


#Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1


#Set master ip to /etc/hosts
echo "$MASTER_NODE  hadoop_master" >> /etc/hosts
