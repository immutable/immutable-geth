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

package params

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
)

// IsImmutableZKEVM returns true if the configured chain ID pertains to an Immutable zkEVM network ID
func (c *ChainConfig) IsImmutableZKEVM() bool {
	return c != nil && c.ChainID != nil &&
		(c.ChainID.Int64() == settings.MainnetNetworkID ||
			c.ChainID.Int64() == settings.TestnetNetworkID ||
			c.ChainID.Int64() == settings.DevnetNetworkID)
}

// IsImmutableZKEVMPrevrandao returns true if the chain configuration has prevrandao fork enabled
// for the specified timestamp. Mainnet is the only network that has a separate prevrandao fork
// to the shanghai fork.
func (c *ChainConfig) IsImmutableZKEVMPrevrandao(blockTime uint64) bool {
	if !c.IsImmutableZKEVM() {
		return false
	}
	if c.PrevrandaoTime == nil {
		panic(fmt.Sprintf("prevrandaoTime not set for Immutable zkEVM network %d", c.ChainID.Int64()))
	}
	return *c.PrevrandaoTime <= blockTime
}

// IsImmutableZKEVMShanghai returns true if the chain configuration has shanghai fork enabled
// and the specified block number and timestamp is past the shanghai timestamp
func (c *ChainConfig) IsImmutableZKEVMShanghai(blockNum *big.Int, blockTime uint64) bool {
	return c.IsImmutableZKEVM() && c.IsShanghai(blockNum, blockTime)
}

// IsValidImmutableZKEVM returns true if the chain configuration is valid for an Immutable zkEVM network
func (c *ChainConfig) IsValidImmutableZKEVM() bool {
	return c.IsImmutableZKEVM() &&
		c.IsReorgBlocked &&
		c.Clique != nil &&
		c.Clique.Period == settings.SecondsPerBlock
}
