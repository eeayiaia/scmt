#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

send_metadata_interval=30

globals_attr="
  daemonize = yes
  setuid = yes
  user = nobody
  debug_level = 0
  max_udp_msg_len = 1472
  mute = no
  deaf = no
  allow_extra_data = yes
  host_dmax = 86400 /*secs. Expires (removes from web interface) hosts in 1 day$
  host_tmax = 20 /*secs */
  cleanup_threshold = 300 /*secs */
  gexec = no
  # By default gmond will use reverse DNS resolution when displaying your hostn$
  # Uncommeting following value will override that value.
  # override_hostname = \"mywebserver.domain.com\"
  # If you are not using multicast this value should be set to something other $
  # Otherwise if you restart aggregator gmond you will get empty graphs. 60 sec$
  send_metadata_interval = $send_metadata_interval /*secs */
"

cluster_attr="
  name = ${1-my cluster}
  owner = ${2-unspecified}
  latlong = ${3-unspecified}
  url = ${4-unspecified}
"
#Assumes accessible them ip for eth0 is accessible for master.
#Change to localhost possible?
udp_send_channel_attr="
  host = localhost
  port = 8649
  ttl = 1
"

udp_recv_channel_attr="
  port = 8649
"

tcp_accept_channel_attr="
  port = 8649
"

data_source="data_source \"${1-my cluster}\" localhost\n"

check_root

apt-get install -y apache2 rrdtool ganglia-monitor ganglia-monitor-python gmetad ganglia-webfrontend    


INSTALL_SUCCESS=$?

if [[ $INSTALL_SUCCESS != 0 ]]; then
    echo "Failed to install ganglia master." >&2
    exit 1
fi

python helpscript/regex.py "gmond" "globals" "$globals_attr"
python helpscript/regex.py "gmond" "cluster" "$cluster_attr"
python helpscript/regex.py "gmond" "udp_send_channel" "$udp_send_channel"
python helpscript/regex.py "gmond" "udp_recv_channel" "$udp_recv_channel_attr"
python helpscript/regex.py "gmond" "tcp_accept_channel" "$tcp_accept_channel_attr"
python helpscript/regex.py "gmetad" "$data_source"

ln -s -f /etc/ganglia-webfrontend/apache.conf /etc/apache2/conf-enabled/ganglia.conf

service ganglia-monitor apache2 gmetad restart