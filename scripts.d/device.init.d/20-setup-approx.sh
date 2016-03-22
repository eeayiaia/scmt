#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../plugins.d/.script-utils/installer-utils.sh" || exit 1

check_root

SOURCES_FILE="$DIR/../../config/approx/node-sources"
SOURCES_TARGET="/etc/apt/sources.list"
SOURCES_D_TARGET="/etc/apt/sources.list.d"

if [[ ! -f $SOURCES_FILE ]]; then
	echo "Could not find '<scmt-root>/config/approx/node-sources', aborting approx setup"
	exit 2
fi

backup_file $SOURCES_TARGET
OLD_SOURCES_LIST=$BACKUP_OUTPUT
cp -f $SOURCES_FILE $SOURCES_TARGET

backup_file $SOURCES_D_TARGET
OLD_SOURCES_LIST_D=$BACKUP_OUTPUT
rm -rf $SOURCES_D_TARGET

# apt-get update returns 0 even when it failed... workaround with grep
UPDATE_SUCCESS=$(apt-get update | grep "Err")

if [[ $UPDATE_SUCCESS != 0 ]]; then
	echo "apt-get update failed after approx install, please make sure your approx configuration '<scmt-root>/config/approx/' is set up properly for your cluster."
	echo ""
	echo "Reverting to old sources.list..."

	cp -f $OLD_SOURCES_LIST $SOURCES_TARGET
	cp -rf $OLD_SOURCES_LIST_D $SOURCES_D_TARGET

	exit 3
fi

