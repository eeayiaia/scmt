#!/bin/bash

# Input: MYSQL_PASSWORD

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

echo "Installing mysql..."
write_line

echo "mysql-server mysql-server/root_password password $MYSQL_PASSWORD" \
	| debconf-set-selections
echo "mysql-server mysql-server/root_password_again password $MYSQL_PASSWORD" \
	| debconf-set-selections

echo "Creating SCMT database..."

mysql -u root -p$MYSQL_PASSWORD < "$DIR/resources/create_database.sql"

