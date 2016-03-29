#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

nodeIP=$1

sed -i "/^data_source/ s/$/ $nodeIP/" /etc/ganglia/gmetad.conf