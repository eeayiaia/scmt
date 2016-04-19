#!/bin/bash

# Get script directory & include utils
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../utils.sh" || exit 1

check_invoked_by_scmt

function before_reboot() {
	# Expands the root partition table to its maximum size
	start=`fdisk -l /dev/mmcblk0 | grep mmcblk0p2 | awk '{print $2}'`
	echo "Found the start of mmcblk0p2: $start"

fdisk /dev/mmcblk0 << __EOF__ >> /dev/null
d
2
n
p
2
$start

p
w
__EOF__

	sync
	touch /root/.resize
	echo "Partition table successfully resized"
}

function after_reboot() {
	echo "Expanding filesystem!"
	write_line
	resize2fs /dev/mmcblk0p2 >> /dev/null
	write_line
	echo "Expanded filesystem!"
}

if [[ ! -f /root/.resize ]]; then
	before_reboot
	reboot now
else
	after_reboot
fi

