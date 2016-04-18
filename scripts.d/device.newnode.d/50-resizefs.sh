#!/bin/sh

[![ -f /root/.resize ]] && exit 1

echo "Expanding filesystem!"
resize2fs /dev/mmcblk0p2 >> /dev/null
echo "Expanded filesystem!"
