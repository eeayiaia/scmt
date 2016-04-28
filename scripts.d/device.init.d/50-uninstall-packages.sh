#!/bin/bash

# Get script directory & include utils
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

REMOVE_LIST_FILE=$CONFIG_DIR/package-remove-list.txt

echo "Uninstalling packages listed in $REMOVE_LIST_FILE..."

if [[ ! -f "$REMOVE_LIST_FILE" ]]; then
	echo "Error: could not find $REMOVE_LIST_FILE" >&2
	exit 1
fi

# Hack to build package array from config file
echo "REMOVE_LIST=(" > temp
cat "$REMOVE_LIST_FILE" >> temp
echo ")" >> temp

source temp
delete_file temp

# Uninstall packages
for i in "${REMOVE_LIST[@]}"; do
	write_line
	echo "Uninstalling $i"
	write_line

	apt-get purge "$i" --assume-yes
done

echo "Finished uninstalling packages."

