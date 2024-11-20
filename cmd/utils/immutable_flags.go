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

package utils

import (
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/ethereum/go-ethereum/triedb/hashdb"
	"github.com/ethereum/go-ethereum/triedb/pathdb"
)

// MakeTrieDatabaseNoContext is MakeTrieDatabase without being coupled to cli context to allow for further reuse.
func MakeTrieDatabaseNoContext(stateScheme string, disk ethdb.Database, preimage bool, readOnly bool, isVerkle bool) *triedb.Database {
	config := &triedb.Config{
		Preimages: preimage,
		IsVerkle:  isVerkle,
	}
	scheme, err := rawdb.ParseStateScheme(stateScheme, disk)
	if err != nil {
		Fatalf("%v", err)
	}
	if scheme == rawdb.HashScheme {
		// Read-only mode is not implemented in hash mode,
		// ignore the parameter silently. TODO(rjl493456442)
		// please config it if read mode is implemented.
		config.HashDB = hashdb.Defaults
		return triedb.NewDatabase(disk, config)
	}
	if readOnly {
		config.PathDB = pathdb.ReadOnly
	} else {
		config.PathDB = pathdb.Defaults
	}
	return triedb.NewDatabase(disk, config)
}
