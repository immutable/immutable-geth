#!/bin/bash

# These tests bootstrap a local network and run tests through the eth JSON RPC.

set -e
set -o pipefail

# necho prints a line without a newline
function necho() {
  echo -n "> $*: "
}

# Clean up subprocesses on exit
function stop_geth() {
  pkill geth || true
  # Wait for kill
  sleep 1
  while [ "$(pgrep geth | wc -l)" -gt 1 ]; do
    sleep 1
    echo -n "."
  done
}

function _exit() {
  stop_geth
}
trap _exit EXIT

# Assign in case PWD changes
dir="$PWD"
log="/tmp/rewindout"
boots=0
validators=1
rpcs=0
now=$(date +%s)
prevrandao_timestamp=$now
shanghai_timestamp=$now
cancun_timestamp=$now
root_dir="$(mktemp -d)"

function start_geth() {
  necho "Starting geth"
  dir="$root_dir/devnet/chain-15003/validator-0"
  export GETH_FLAG_PASSWORD_FILEPATH="$dir/password"

  ./build/bin/geth \
  --datadir "$dir" --log.debug --networkid 15003 --metrics --metrics.addr 127.0.0.1 \
  --metrics.port 6060 --authrpc.port 8550 --verbosity 4 --port 30300 \
  --rpc.debugdisable --rpc.txpooldisable \
  --config "$root_dir/devnet/chain-15003/config.toml" --pprof --pprof.port 7070 \
  --miner.etherbase "$(cat "$dir/address")" \
  --mine --http --http.port 8545 --cache 128 --cache.database 35 \
  --cache.trie 35 --cache.gc 10 --cache.snapshot 20 --gcmode archive --syncmode full \
  --override.prevrandao="$prevrandao_timestamp" \
  --override.shanghai="$shanghai_timestamp" \
  --override.cancun="$cancun_timestamp" >> "$log" 2>&1 &

  # Wait for geth processes to start
  while [ "$(pgrep geth | wc -l)" -lt "$((boots+rpcs+1))" ]; do
    sleep 1
    echo -n "."
  done
  echo ""
  necho "To view geth logs run"
  echo "tail -f $log"
}

function bootstrap_geth() {
  necho "Bootstrapping geth"

  # Bootstrap and run local network
  ./build/bin/geth immutable bootstrap local \
  --datadir "$root_dir" \
  --boots "$boots" \
  --rpcs "$rpcs" \
  --validators "$validators" \
  --blocklistfilepath "$dir/cmd/geth/testdata/acl_list.txt" \
  --override.prevrandao="$prevrandao_timestamp" \
  --override.shanghai="$shanghai_timestamp" \
  --override.cancun="$cancun_timestamp" >> "$log" 2>&1 &

  # Wait for geth processes to start
  while [ "$(pgrep geth | wc -l)" -lt "$((boots+rpcs+1))" ]; do
    sleep 1
    echo -n "."
  done
  echo ""
  necho "To view geth logs run"
  echo "tail -f $log"
}

function get_head() {
  resp=$(curl -Ss --location 'http://127.0.0.1:8545' \
  --header 'Content-Type: application/json' \
  --data '{
  	"jsonrpc":"2.0",
  	"method":"eth_getBlockByNumber",
  	"params":[
  		"latest", 
  		true
  	],
  	"id":1
  }')
  block=$(echo "$resp" | jq -r '.result.number')
  block=$((block))
  if [ "$block" == 0 ]; then
    echo "Chain not progressing, exiting"
    exit 1
  fi
}

# Run geth
rm "$log" || true
bootstrap_geth

# Check block after a few should have been sealed
sleep 20
get_head
pre_rewind_block="$block"
echo "Block before rewind: $pre_rewind_block"

# Kill the geth process and rewind
echo "Killing geth"
stop_geth

echo "Rewinding chain"
rewind_block=0
./build/bin/geth immutable rewind \
--datadir "$root_dir/devnet/chain-15003/validator-0" \
"$rewind_block" >> "$log" 2>&1
cat "$root_dir/devnet/chain-15003/validator-0/rewind_history.yaml"

# Restart geth and wait for smaller period of time than pre rewind
start_geth
sleep 8

# Check latest header and compare against pre-rewind block
get_head
echo "Block after rewind: $block"
if [ "$block" -ge "$pre_rewind_block" ]; then
  echo "Latest ($block) is not less than pre-rewind block ($pre_rewind_block)"
  exit 1
fi

# Kill the geth process and rewind to same block again, expecting no change
sleep 10 # Wait for chain to progress further beyond `block`
echo "Killing geth"
stop_geth

echo "Rewinding chain again"
pre_rewind_block="$block"
rewind_block=0
./build/bin/geth immutable rewind \
--datadir "$root_dir/devnet/chain-15003/validator-0" \
"$rewind_block" >> "$log" 2>&1
cat "$root_dir/devnet/chain-15003/validator-0/rewind_history.yaml"

start_geth
sleep 8
get_head
echo "Block after second rewind: $block"
if [ "$block" -le "$pre_rewind_block" ]; then
  echo "Latest ($block) should be greater than pre rewind block ($pre_rewind_block)"
  exit 1
fi

echo "Tests finished successfully"
