#!/bin/bash
set -e

NODE_NAME=simappcli
RELAYER_CMD="${RELAYER_CLI} --home $(pwd)/.relayer"
CONFIG_PATH=$(pwd)/../configs/demo

# First initialize your configuration for the relayer
${RELAYER_CMD} config init

# Then add the chains and paths that you will need to work with the
# gaia chains spun up by the two-chains script
${RELAYER_CMD} chains add -f ${CONFIG_PATH}/ibc0.json
${RELAYER_CMD} chains add -f ${CONFIG_PATH}/ibc1.json

# To finalize your config, add a path between the two chains
${RELAYER_CMD} paths add ibc0 ibc1 demo -f ${CONFIG_PATH}/demo.json

# Now, add the key seeds from each chain to the relayer to give it funds to work with
${RELAYER_CMD} keys restore ibc0 testkey "$(jq -r '.secret' data/ibc0/n0/${NODE_NAME}/key_seed.json)"
${RELAYER_CMD} keys restore ibc1 testkey "$(jq -r '.secret' data/ibc1/n0/${NODE_NAME}/key_seed.json)"

# Then its time to initialize the relayer's lite clients for each chain
# All data moving forward is validated by these lite clients.
${RELAYER_CMD} lite init ibc0 -f
${RELAYER_CMD} lite init ibc1 -f

${RELAYER_CMD} tx link demo

${RELAYER_CMD} start demo