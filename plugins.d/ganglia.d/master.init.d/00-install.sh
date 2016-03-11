#!/bin/bash

cluster_attr = "cluster{
  name = ${2-my cluster}
  owner = ${3-unspecified}
  latlong = ${4-unspecified}
  url = ${5-unspecified}
}"

udp_send_channel_attr = "udp_send_channel { 
  host =  $(awk '/^[[:space:]]*($|#)/{next} /master/{print $1; exit}' /etc/hosts)
  port = 8649
  ttl = 1
} "


# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

check_root

apt-get install ganglia-monitor ganglia-monitor-python gmetad


INSTALL_SUCCESS=$?

if [[ $INSTALL_SUCCESS != 0 ]]; then
    echo "Failed to install Ganglia." >&2
    exit 1
fi

