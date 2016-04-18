#!/bin/sh

# Stash current changes to prevent tests
# being run on experimental code
git stash -q --keep-index

go test ./...
RESULT=$?

git stash pop -q

[ $RESULT -ne 0 ] && exit 1
exit 0
