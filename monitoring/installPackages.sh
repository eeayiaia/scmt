#!/bin/sh

if [ "$#" -le 1 ]; then
    echo "Usage: ./$0 <package directory> <packagename1> <packagename2> .... <packagenameN>"
fi

installDirectory=$1

for i in ${@:2}
do
    dpkg -i "$0$i"
done
