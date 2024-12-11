#!/bin/bash

set -e
set -o pipefail

# Input args
env="$1"
if [ -z "$env" ]; then
  echo "Environment (dev, sandbox) must be provided as arg"
  exit 1
fi
release_tag="$2"
if [ -z "$release_tag" ]; then
  echo "Release tag must be provided as arg"
  exit 1
fi
echo "Checking for release tag ($release_tag)"

# Main loop
tag=""
loop_sleep=10
while [ "$tag" != "$release_tag" ]; do
  # Get all images pertaining to zkevm-geth. Sort them and get unique values.
  image_list=$(kubectl get pods -n "$env" -lapp=zkevm-geth -oyaml | grep image: | awk '{print $2}' | sort -u | grep go-ethereum)

  # No geth images found
  if [ -z "$image_list" ]; then
    echo "No go-ethereum images found"
    exit 1
  fi

  # More than one geth image found
  if [ "$(echo "$image_list" | wc -l)" -gt 1 ]; then
    echo "Multiple go-ethereum images found, rollout in progress: "
    echo "$image_list"
    sleep "$loop_sleep"
    continue
  fi

  # We now know that it is a single image, extract the tag
  tag=$(echo "$image_list" | sed 's/.*://')

  # Check if the release tag matches the image tag
  if [ "$tag" != "$release_tag" ]; then
    echo "Release tag ($release_tag) does not match image tag ($tag)"
    sleep "$loop_sleep"
    continue
  fi
done

echo "Release tag matches image tag ($tag)"
