#!/bin/bash

# Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

# Install aptitude
write_line
echo "Installing aptitude..."
apt-get install aptitude --assume-yes
[[ ! $? ]] && echo "Failed to install aptitude" >&2
write_line

# Install realpath
write_line
echo "Installing realpath..."
apt-get install realpath --assume-yes
[[ ! $? ]] && echo "Failed to install realpath" >&2
write_line

# Install mysql-server
write_line
echo "mysql-server mysql-server/root_password password $MYSQL_ROOT_PASSWORD" | sudo debconf-set-selections
echo "mysql-server mysql-server/root_password_again password $MYSQL_ROOT_PASSWORD" | sudo debconf-set-selections
echo "Installing mysql-server..."
apt-get install mysql-server --assume-yes
[[ ! $? ]] && echo "Failed to install mysql-server" >&2
write_line

