#!/bin/bash
set -e

PROCESS_NAME=simappd

rm -rf ./.relayer
rm -rf ./data

count=`ps -ef | grep $PROCESS_NAME | grep -v grep | wc -l`

if [ $count != 0 ]; then
    echo "killall $PROCESS_NAME."
    killall ${PROCESS_NAME}
else
    echo "$PROCESS_NAME not yet run."
fi
