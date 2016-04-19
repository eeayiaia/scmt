#!/bin/bash

find ./ \( -name "*.go" -or -name "*sh" -and -not -name "detect-todo.sh" \) \
	| xargs egrep -i "(\W)(TODO|FIX|HACK)"

