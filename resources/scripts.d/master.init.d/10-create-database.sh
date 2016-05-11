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

MYSQL_QUERY="CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE;
USE $MYSQL_DATABASE;"'
CREATE TABLE IF NOT EXISTS devices (hwaddr CHAR(12) KEY, id INT UNSIGNED AUTO_INCREMENT, port INT UNSIGNED, hname varchar(30), username VARCHAR(20), password VARCHAR(20), KEY `id` (`id`));
CREATE TABLE IF NOT EXISTS plugins (name VARCHAR(30) PRIMARY KEY,    enabled BOOLEAN NOT NULL DEFAULT 0, installedOnMaster BOOLEAN NOT NULL DEFAULT 0);
CREATE TABLE IF NOT EXISTS installedPlugins_slave (hwaddr CHAR(12) NOT NULL, plugin VARCHAR(30) NOT NULL, FOREIGN KEY (hwaddr) REFERENCES devices(hwaddr), FOREIGN KEY (plugin) REFERENCES plugins(name), UNIQUE (hwaddr, plugin));'"
GRANT ALL PRIVILEGES ON $MYSQL_DATABASE . * TO $MYSQL_USER@localhost IDENTIFIED BY '$MYSQL_PASSWORD';
FLUSH PRIVILEGES;"

echo "Creating SCMT database..."

echo "$MYSQL_QUERY" | mysql -u root -p$MYSQL_ROOT_PASSWORD

