#!/bin/bash

set -e

# Validation relevant to all environments
for env in dev sandbox prod; do
  # Check that the validator replica count is 1
  replicas="$(grep "replicas:" < "./deployment/$env/geth/validator/statefulset.yaml" | sed "s/.*: //g")"
  # In the case of validator rotation, we need to disable this invariant check
  # Revert back to 1 after the rotation is complete
  if [ "$replicas" != "1" ]; then
    echo "$env: Validator replica count is not empty or set to 1"
    exit 0
  fi

  # Check config.toml
  max=999999999000000000
  conf="./deployment/$env/geth/configmaps/config.toml"
  if ! grep -q "Recommit = $max" < "$conf"; then
    echo "$env: Recommit is not set to $max"
    exit 1
  fi
  price=10000000000
  if ! grep -q "GasPrice = $price" < "$conf"; then
    echo "$env: GasPrice is not set to $price"
    exit 1
  fi
  if ! grep -q "PriceLimit = $price" < "$conf"; then
    echo "$env: PriceLimit is not set $price"
    exit 1
  fi
done

# Check csvs have valid addresses
for csv in sdn_list.txt blocklist_manual_list.txt; do
  # Eth addr regex
  regex='^0x([a-fA-F0-9]){40}$'

  # Read csv into array
  IFS="," read -r -a addrs <<< "$(cat "./deployment/base/geth/configmaps/$csv")"

  # Validate the addrs in the array
  len=${#addrs[@]}
  if [ "$len" -eq 0 ]; then
    echo "$env: No addresses found in $csv"
    exit 1
  fi
  for addr in "${addrs[@]}"; do
    if [[ ! "$addr" =~ $regex ]]; then
      echo "$env: invalid address: $addr"
      exit 1
    fi
  done
  echo "$env: Addresses ($len) in $csv are valid"
done

# Check version of geth src is aligned with changelog
changelog_version=$(grep '## \[' < CHANGELOG.md | head -n1 | tr -d ' #[]')
major=$(grep '// Major version component' < params/version.go | sed "s|//.*||" | sed "s|.* = ||")
minor=$(grep '// Minor version component' < params/version.go | sed "s|//.*||" | sed "s|.* = ||")
patch=$(grep '// Patch version component' < params/version.go | sed "s|//.*||" | sed "s|.* = ||")
meta=$(grep '// Version metadata to append' < params/version.go | sed "s|//.*||" | sed "s|.* = ||")
src_version=$(echo "v$major.$minor.$patch-$meta" | tr -d '\t" ')
if [ "$src_version" != "$changelog_version" ]; then
  echo "Version in params/version.go ($src_version) is not aligned with CHANGELOG.md ($changelog_version)"
  exit 1
fi
echo "Version in params/version.go is aligned with CHANGELOG.md ($changelog_version)"
