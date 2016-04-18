#!/bin/sh

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
