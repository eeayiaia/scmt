#!/bin/bash

# Get script directory and follow symlinks
DIR=$(dirname $(readlink -f $0))
export SCMT_ROOT="$DIR/resources"
$DIR/scmt "$@"
