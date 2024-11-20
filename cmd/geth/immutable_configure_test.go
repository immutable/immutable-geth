// Copyright 2024 The Immutable go-ethereum Authors
// This file is part of the Immutable go-ethereum library.
//
// The Immutable go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Immutable go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Immutable go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/role"
)

func TestImmutableConfigure(t *testing.T) {
	// Skip this test if rotating the validator keys. Invariants will be disabled.
	// t.Skip()
	dir := filepath.Join(os.TempDir(), "immutable-test")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		name      string
		r         role.Role
		pod       string
		ordinal   int
		filepaths configureNodePaths
		err       error
	}{
		{
			name:    "validator valid config",
			r:       role.Validator,
			pod:     "zkevm-geth-validator-0",
			ordinal: 0,
			filepaths: configureNodePaths{
				ConfigFilepath: "../../cmd/geth/testdata/dev.toml",
				DataDirpath:    dir,
			},
			err: nil,
		},
		{
			name:    "validator invalid config",
			r:       role.Validator,
			pod:     "zkevm-geth-validator-0",
			ordinal: 0,
			filepaths: configureNodePaths{
				ConfigFilepath: "../../cmd/geth/testdata/dev_bad_recommit.toml",
				DataDirpath:    dir,
			},
			err: ErrConfigRecommitInvariant,
		},
		{
			name:    "validator invalid ordinal",
			r:       role.Validator,
			pod:     "zkevm-geth-validator-1",
			ordinal: 1,
			filepaths: configureNodePaths{
				ConfigFilepath: "../../cmd/geth/testdata/dev.toml",
				DataDirpath:    dir,
			},
			err: ErrOrdinalInvariant,
		},
		{
			name:    "boot valid config",
			r:       role.Boot,
			pod:     "zkevm-geth-boot-0",
			ordinal: 0,
			filepaths: configureNodePaths{
				ConfigFilepath: "../../cmd/geth/testdata/dev.toml",
				DataDirpath:    dir,
			},
			err: nil,
		},
		{
			name:    "boot invalid ordinal",
			r:       role.Boot,
			pod:     "zkevm-geth-boot-1",
			ordinal: 1,
			filepaths: configureNodePaths{
				ConfigFilepath: "../../cmd/geth/testdata/dev.toml",
				DataDirpath:    dir,
			},
			err: ErrOrdinalInvariant,
		},
		{
			name:    "rpc valid config",
			r:       role.RPC,
			pod:     "zkevm-geth-rpc-2",
			ordinal: 2,
			filepaths: configureNodePaths{
				ConfigFilepath: "../../cmd/geth/testdata/dev.toml",
				DataDirpath:    dir,
			},
			err: nil,
		},
		{
			name:    "partner valid config",
			r:       role.Partner,
			pod:     "zkevm-geth-partner-3",
			ordinal: 3,
			filepaths: configureNodePaths{
				ConfigFilepath: "../../cmd/geth/testdata/dev.toml",
				DataDirpath:    dir,
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := configureNodeInPod(test.r, test.pod, test.filepaths)
			if err != nil {
				if test.err == nil {
					t.Fatal(err)
				}
				if !errors.Is(err, test.err) {
					t.Fatalf("error: %v, expected: %v", err, test.err)
				}
				// Got the error we expected
				return
			} else if test.err != nil {
				t.Fatalf("expected error: %v", test.err)
			}
		})
	}
}
