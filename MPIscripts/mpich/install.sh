#!/bin/bash

clear

echo "Do you wish to install MPICH? [y/n]"

read answer

case $answer in
	[Yy]* ) sh ./installMPICH.sh;;
	[Nn]* ) echo "MPICH install aborted."; break;;
esac

echo " "
echo "returning to previus meny." 
