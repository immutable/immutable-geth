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

package core

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
)

// ImmutableGenesisBlock returns the immutable genesis block for the specified network (mainnet, testnet, devnet).
func ImmutableGenesisBlock(network string) *Genesis {
	net, err := settings.NewNetwork(network)
	if err != nil {
		panic(err)
	}
	// Unmarshal the JSON string into a Genesis struct
	genesis := &Genesis{}
	if err := json.Unmarshal([]byte(net.GenesisJSON()), genesis); err != nil {
		panic(err)
	}
	return genesis
}
