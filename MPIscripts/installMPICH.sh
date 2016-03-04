#!/bin/bash

whereis gcc | wc -l > 0 |  echo "whoo!"
echo 
if [ $(whereis gcc | wc -m) -ge 6 ]
then
	echo "gcc is installed"
else
	echo "gcc is not installed do it.."
fi
