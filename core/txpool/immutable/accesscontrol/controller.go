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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Controller struct {
	// Map of filename to address provider
	providers map[string]AddressProvider
	// isAnAllowList indicates that it is an allowlist otherwise the inverse is a blocklist
	isAnAllowList bool
}

// New initializes an access controller with blocklists specified by file paths.
//
// Parameters:
//   - filePaths: A slice of strings containing file paths to blocklist
//     files, usually an sdn file that comes in the format of txt as comma separated values
//   - isAnAllowList: Indicates if the controller is an allow controller or
//     a block controller.
func New(filePaths []string, isAnAllowList bool) (*Controller, error) {
	providers := make(map[string]AddressProvider, len(filePaths))

	for _, filename := range filePaths {
		csvProvider, err := newCSVProvider(filename)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialize access controller provider: %w", err)
		}
		providers[filename] = csvProvider
	}

	return &Controller{
		providers:     providers,
		isAnAllowList: isAnAllowList,
	}, nil
}

func (c *Controller) IsBlocklist() bool {
	return !c.isAnAllowList
}

func (c *Controller) IsAllowed(addr common.Address, tx *types.Transaction) bool {
	for _, list := range c.providers {
		addresses := list.Provide()
		if _, exist := addresses[addr]; exist {
			return c.isAnAllowList
		}
		if c.IsBlocklist() && tx.To() != nil {
			if _, exist := addresses[*tx.To()]; exist {
				return c.isAnAllowList
			}
		}
	}

	// If the address is not in the list and it's not an allow list, return true
	return !c.isAnAllowList
}
