#!/bin/bash

set -e

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
