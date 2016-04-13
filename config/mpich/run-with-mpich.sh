#!/bin/bash

if [[ ! $1 ]]; then
	echo "Usage: run-with-mpich <executable> <no. of processes>"
	exit 1
fi

mpirun.mpich -f ~/mpich-machinefile -n $2 ./$1

