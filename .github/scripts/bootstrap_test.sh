#!/bin/bash

# These tests bootstrap a local network and run tests through the eth JSON RPC.
# The network will undergo a number of hard forks between various test fixtures.
# The Prevrandao fork is intended to be <= Shanghai fork.

# TODO(serge): move these tests to a proper test suite implemented in Go
# when we know what the final test suite will look like.
# Features for test suite:
# * Set of fixtures that create transactions affected by all EVM-related hard forks
# * Fixtures can be run with expectation of success or failure depending on fork state

set -e
set -o pipefail

# necho prints a line without a newline
function necho() {
  echo -n "> $*: "
}

# Clean up subprocesses on exit
_exit() {
  pkill geth || true
  # sleep 5s to let it kill
  sleep 5
}
trap _exit EXIT

# Assign in case PWD changes
dir="$PWD"
log="/tmp/bootstrapout"
boots=2
validators=1
rpcs=2
# Set prevrandao override to be before shanghai
now=$(date +%s)
prevrandao_timestamp=$now
# Set shanghai override to be after genesis block and pre-shanghai tests
shanghai_timestamp=$((now+30))
# Set cancun override a few blocks after shanghai
cancun_timestamp=$((shanghai_timestamp+4))

export GETH_FLAG_NET_RESTRICT="0.0.0.0/0"
export GETH_FLAG_P2P_SUBNET="127.0.0.1/32"

function start_geth() {
  necho "Starting geth"

  # Bootstrap and run local network
  ./build/bin/geth immutable bootstrap local \
  --boots "$boots" \
  --rpcs "$rpcs" \
  --validators "$validators" \
  --blocklistfilepath "$dir/cmd/geth/testdata/acl_list.txt" \
  --override.prevrandao="$prevrandao_timestamp" \
  --override.shanghai="$shanghai_timestamp" \
  --override.cancun="$cancun_timestamp" > "$log" 2>&1 &

  # Wait for geth processes to start
  while [ "$(pgrep geth | wc -l)" -lt "$((boots+rpcs+1))" ]; do
    sleep 1
    echo -n "."
  done
  echo ""
  necho "To view geth logs run"
  echo "tail -f $log"
  sleep 2
}

function wait_shanghai() {
  # Wait for shanghai timestamp
  now=$(date +%s)
  if [ "$now" -lt "$shanghai_timestamp" ]; then
    diff=$((shanghai_timestamp-now+5))
    necho "Waiting $diff seconds for shanghai timestamp"
    sleep "$diff"
  fi
  echo ""
}

# Run geth
start_geth

# Need to use 0x44 opcode for Prevrandao fork.
echo "> Running post-prevrandao, pre-shanghai tests"
go test -count=1 -v ./tests/immutable \
-privkey="$dir/cmd/geth/testdata/key.prv" \
-rpc=http://localhost:8546 \
-validatoradmin=http://localhost:8545 \
-skipvoting=true \
-run='.*Randao.*'

# Wait for shanghai override
wait_shanghai

# Run tests
# Need to use contracts compiled with solc relevant to Shanghai fork.
echo "> Running post-fork tests"
go test -count=1 -v ./tests/immutable \
-privkey="$dir/cmd/geth/testdata/key.prv" \
-blockedprivkey="$dir/cmd/geth/testdata/blockedkey.prv" \
-rpc=http://localhost:8546 \
-validatoradmin=http://localhost:8545 \
-skipvoting=true

echo "> Running vote test" # Tests after this will fail due to stalled block production
./build/bin/geth immutable vote add --voters http://localhost:8545 --validator 0x7442eD1e3c9FD421F47d12A2742AfF5DaFBf43f8
echo "Tests finished successfully"
