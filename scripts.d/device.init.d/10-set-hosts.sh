#!/bin/bash

# Input: MASTER_IP
# Set the hosts file to include the master node

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

egrep -q "^master\s" /etc/hosts \
	&& sed "s/master/$MASTER_IP    master/" -i /etc/hosts \
	|| sed "$ a\$MASTER_IP    master" -i /etc/hosts

