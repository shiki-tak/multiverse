#!/bin/bash
set -e

PREV_DIR=$(pwd)
RELAYER_DIR=$(mktemp -d)

echo "RELAYER_DIR is ${RELAYER_DIR}"

cd ${RELAYER_DIR}
git clone https://github.com/shiki-tak/relayer
cd ./relayer
git checkout connect
echo "Building Relayer..."
make build

export RELAYER_CLI=${RELAYER_DIR}/relayer/build/rly

cd ${PREV_DIR}/scripts
./cleanup.sh

./two-chainz
# wait for all chains to start.
sleep 10
./setup-channel.sh
