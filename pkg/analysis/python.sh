# !/bin/bash

FILENAME=$1
BASENAME=$(basename $FILENAME)

result=$(pylint $FILENAME | grep $BASENAME | awk -F '/' '{print $NF}')

echo "$result"
