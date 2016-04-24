#!/bin/bash

# Input: MYSQL_PASSWORD

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

mysql -u root -p$MYSQL_PASSWORD < "$DIR/resources/create_database.sql"

