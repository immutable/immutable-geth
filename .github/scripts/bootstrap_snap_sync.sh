#!/bin/bash

# This script can be used to start a local geth node from genesis against a specified
# network. For example, `./.github/scripts/dev.toml` is a config for connecting to devnet
# if you have port-forwarded to the p2p partner pod via `kubectl port-forward pod/zkevm-geth-partner-0 30300:30300`.

set -e
set -o pipefail

# Clean up subprocesses on exit
_exit() {
  pkill geth || true
  # sleep 5s to let it kill
  sleep 5
}
trap _exit EXIT

# Set the env var to enable long range sync
export GETH_FLAG_IMMUTABLE_LONG_RANGE_SYNC="1"

./build/bin/geth immutable bootstrap local \
--env devnet \
--syncmode snap \
--gcmode full \
--config ./.github/scripts/dev.toml \
--boots "0" \
--rpcs "1" \
--validators "0"