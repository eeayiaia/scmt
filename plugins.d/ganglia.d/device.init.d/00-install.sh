
# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../.script-utils/installer-utils.sh"

cluster_attr="
  name = ${2-my cluster}
  owner = ${3-unspecified}
  latlong = ${4-unspecified}
  url = ${5-unspecified}
"

udp_send_channel_attr="
  host =  $(awk '/^[[:space:]]*($|#)/{next} /master/{print $1; exit}' /etc/hosts)
  port = 8649
  ttl = 1
"

check_root

apt-get install -y ganglia-monitor

INSTALL_SUCCESS=$?

if [[ $INSTALL_SUCCESS != 0 ]]; then
    echo "Failed to install ganglia-monitor." >&2
    exit 1
fi


python helpscript/regex.py "$cluster_attr" "$udp_send_channel_attr"

service ganglia-monitor restart

#sed -i 's/clustr{/new-word/g' /etc/ganglia/gmond.conf
