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

apt-get install mysql-server --assume-yes
INSTALL_SUCCESS=$?

if [[ ! $INSTALL_SUCCESS ]]; then
	echo "Failed to install mysql." >&2
	exit 1
fi

echo "Creating SCMT database..."

mysql -u root -p$MYSQL_PASSWORD < "$DIR/resources/create_database.sql"

