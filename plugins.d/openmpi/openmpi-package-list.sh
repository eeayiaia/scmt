#!/bin/bash

# This package list is everything needed to install OpenMPI on a freshly
# installed desktop Ubuntu 15.10.
# TODO: Generalise somehow?

# Note that package must be listed in order of installation - some are
# dependent on others.
packages=(libcr0 libibverbs1 libibverbs-dev libnuma-dev libhwloc5 libhwloc-dev libopenmpi1.6 openmpi-common openmpi-bin libopenmpi-dev)

