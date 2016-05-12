#!/bin/bash

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

INITFILE="resources/scmt"
INITDPATH="/etc/init.d/scmt"

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

SCMT_ROOT="$DIR/../.."

echo "Copying scmt init file to $INITDPATH"
sed -e "s|\$SCMT_ROOT|$SCMT_ROOT| " "$INITFILE" > "$INITDPATH"
chmod 755 $INITDPATH

echo "Update init record"
update-rc.d "$(basename $INITDPATH)" defaults

