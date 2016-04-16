#!/bin/bash

if [[ ! $1 ]]; then
	echo "Usage: run-with-openmpi <executable> <no. of processes>"
	exit 1
fi

mpirun.openmpi -np $2 --hostfile ~/openmpi-machinefile ./$1

