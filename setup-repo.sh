#/bin/sh

CURRENT_DIR=`pwd`

# Link the pre-commit into git-hooks
ln -s $CURRENT_DIR/pre-commit.sh $CURRENT_DIR/.git/hooks/pre-commit
