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

package localstore

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
)

// Store is used to push configuration artefacts to local files
type Store struct {
}

// New returns a new store
func New() *Store {
	return &Store{}
}

// StoreConfigFile writes a configuration file to disk. We use this to
// update k8s configmaps.
func (s *Store) StoreConfigFile(filepath string, content []byte) error {
	log.Info("Writing config file", "filepath", filepath)
	if err := os.WriteFile(filepath, content, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
