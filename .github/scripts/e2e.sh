#!/bin/bash

set -e
set -o pipefail

rpc="$1"
if [ -z "$rpc" ]; then
      echo "RPC url must be provided as arg"
      exit 1
fi
env=""
case "$rpc" in
  "https://rpc.immutable.com")
    env="mainnet"
    ;;
  "https://rpc.testnet.immutable.com")
    env="testnet"
    ;;
  "https://rpc.dev.immutable.com")
    env="devnet"
    ;;
  *)
    echo "Unknown RPC url: $rpc"
    exit 1
    ;;
esac

if [ -z "$PRIV_KEY" ]; then
  echo "PRIV_KEY env var not set"
  exit 1
fi

go test -count=1 -v ./tests/immutable -rpc="$rpc" -skipvoting=true -forks="$env"
