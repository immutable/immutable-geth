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

package accesscontrol

import (
	"errors"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
)

type CSVProvider struct {
	filePath string

	// Contains a set of blockchain addresses, struct{} is used here as value to optimize memory footprint
	addresses map[common.Address]struct{}
}

func newCSVProvider(filePath string) (*CSVProvider, error) {
	addresses, err := load(filePath)
	if err != nil {
		return nil, err
	}

	return &CSVProvider{
		filePath:  filePath,
		addresses: addresses,
	}, nil
}

func load(filePath string) (map[common.Address]struct{}, error) {
	addresses := make(map[common.Address]struct{})

	log.Info("Loading ACL file", "filepath", filePath)
	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	log.Info("Loaded ACL file", "filepath", filePath, "content", string(byteValue))
	// Split the file content by comma to get individual Ethereum addresses
	ethAddresses := strings.Split(string(byteValue), ",")
	for _, ethAddress := range ethAddresses {
		ethAddress = strings.TrimSpace(ethAddress) // Just to be sure there's no leading or trailing whitespace
		if common.IsHexAddress(ethAddress) {
			addresses[common.HexToAddress(ethAddress)] = struct{}{}
		}
	}
	// if we can't parse any address, the file might be empty or corrupted
	if len(addresses) == 0 {
		return nil, errors.New("file is empty or does not contain any valid addresses")
	}

	return addresses, nil
}

func (s *CSVProvider) Provide() map[common.Address]struct{} {
	return s.addresses
}
