#!/bin/bash

# Supplied environment variables:
#   HOSTNAME

[[ $INVOKED_BY_SCMT == 1 ]] || exit 1

BACKUP_FOLDER=~/.scmt-backup
DATE_STAMP=$(date "+%b_%d_%Y_%H:%M:%S")
BACKUP_FILE=/etc/hostname
BACKUP_OUTPUT=$BACKUP_FOLDER/$BACKUP_FILE-$DATE_STAMP

if [[ ! -d "$BACKUP_FOLDER" ]]; then
	mkdir "$BACKUP_FOLDER"
fi

echo "Backing up file $BACKUP_FILE to $BACKUP_OUTPUT..."
cp "$BACKUP_FILE" "$BACKUP_OUTPUT"

# Set correct hostname
echo $NODE_HOSTNAME > /etc/hostname
echo "Set hostname to $NODE_HOSTNAME"

