#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

cluster_attr="
  name = ${1-my cluster}
  owner = ${2-unspecified}
  latlong = ${3-unspecified}
  url = ${4-unspecified}
"

udp_send_channel_attr="
  host =  $(awk '/^[[:space:]]*($|#)/{next} /master/{print $1; exit}' /etc/hosts)
  port = 8649
  ttl = 1
"

data_source="data_source \"${1-my cluster}\" localhost\n"

check_root

apt-get install -y apache2 rrdtool ganglia-monitor ganglia-monitor-python gmetad ganglia-webfrontend    


INSTALL_SUCCESS=$?

if [[ $INSTALL_SUCCESS != 0 ]]; then
    echo "Failed to install ganglia master." >&2
    exit 1
fi

python helpscript/regex.py "$cluster_attr" "$udp_send_channel_attr" "$data_source"

ln -s -f /etc/ganglia-webfrontend/apache.conf /etc/apache2/conf-enabled/ganglia.conf

service ganglia-monitor apache2 restart