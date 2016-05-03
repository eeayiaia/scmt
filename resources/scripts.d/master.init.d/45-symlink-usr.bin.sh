#!/bin/bash

# Get script directory
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

BIN_PATH="/usr/bin/scmt"

ln -s "$DIR/../../../run-scmt.sh" $BIN_PATH
