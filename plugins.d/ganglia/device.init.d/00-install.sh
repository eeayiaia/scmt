#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

#Default send interval is 5 minutes
send_metadata_interval=${1-300}

globals_attr="
  daemonize = yes
  setuid = yes
  user = nobody
  debug_level = 0
  max_udp_msg_len = 1472
  mute = no
  deaf = yes
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
  name = \"my cluster\"
  owner = unspecified
  latlong = unspecified
  url = unspecified
"
#Assumes master hostname is in /etc/hosts
udp_send_channel_attr="
  host = $(awk '/^[[:space:]]*($|#)/{next} /master/{print $1; exit}' /etc/hosts)
  port = 8649
  ttl = 1
"


check_root

apt-get install -y --force-yes ganglia-monitor

INSTALL_SUCCESS=$?

if [[ $INSTALL_SUCCESS != 0 ]]; then
    echo "Failed to install ganglia-monitor." >&2
    exit 1
fi


ln -s /usr/lib/ganglia/* /usr/lib/

#Backup config file
backup_file /etc/ganglia/gmond.conf

python helpscript/regex.py "gmond" "globals" "$globals_attr"
python helpscript/regex.py "gmond" "cluster" "$cluster_attr"
python helpscript/regex.py "gmond" "udp_send_channel" "$udp_send_channel_attr"

service ganglia-monitor restart
