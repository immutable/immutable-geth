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

package rewind

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/yaml.v2"
)

// Record is a single entry in the rewind history.
type Record struct {
	BlockNumber uint64
	BlockHash   common.Hash
	Timestamp   time.Time
}

// History is a list of all records in the rewind history.
type History []Record

// Contains checks if a block is in the rewind history.
func (h History) Contains(hash common.Hash) bool {
	for _, r := range h {
		if r.BlockHash == hash {
			return true
		}
	}
	return false
}

// ReadRewindHistory reads the rewind history from the expected YAML file.
func ReadRewindHistory(filepath string) (history History, err error) {
	fileData, err := os.ReadFile(filepath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("error reading rewind history file: %v", err)
	}
	if errors.Is(err, fs.ErrNotExist) {
		log.Info("No rewind history file found, creating a new one.")
	} else {
		log.Info("Reading existing rewind history file.")
		if err := yaml.Unmarshal(fileData, &history); err != nil {
			return nil, fmt.Errorf("error unmarshalling rewind history YAML: %v", err)
		}
		return history, nil
	}
	return []Record{}, nil
}

// WriteRewindHistory writes the rewind history to the expected YAML file.
func WriteRewindHistory(history History, filepath string) error {
	yamlData, err := yaml.Marshal(&history)
	if err != nil {
		return fmt.Errorf("error marshalling rewind history YAML: %v", err)
	}
	if err := os.WriteFile(filepath, yamlData, 0644); err != nil {
		return fmt.Errorf("error writing rewind history file: %v", err)
	}
	return nil
}
