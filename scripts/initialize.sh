#!/bin/bash
set -e

PREV_DIR=$(pwd)
RELAYER_DIR=$(mktemp -d)

echo "RELAYER_DIR is ${RELAYER_DIR}"

cd ${RELAYER_DIR}
git clone https://github.com/iqlusioninc/relayer
cd ./relayer
echo "Building Relayer..."
make build

export RELAYER_CLI=${RELAYER_DIR}/relayer/build/rly

cd ${PREV_DIR}

./cleanup.sh

./two-chainz
# wait for all chains to start.
sleep 10
./setup-channel.sh
