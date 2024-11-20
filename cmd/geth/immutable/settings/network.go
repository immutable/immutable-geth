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

package settings

import (
	"fmt"

	_ "embed"
)

var (
	//go:embed genesis/mainnet.json
	immutableGenesisMainnetJSON string
	//go:embed genesis/testnet.json
	immutableGenesisTestnetJSON string
	//go:embed genesis/devnet.json
	immutableGenesisDevnetJSON string
)

// Network is the name of an Immutable zkEVM network
type Network struct {
	string
	genesisJSON string
	id          int
	cancun      Fork
}

// String returns the string representation of the network
func (n Network) String() string {
	return n.string
}

// GenesisJSON returns the JSON string for the genesis block of the network
func (n Network) GenesisJSON() string {
	return n.genesisJSON
}

// ID returns the network ID of the chain
func (n Network) ID() int {
	return n.id
}

// Cancun returns the Cancun fork for the network.
func (n Network) Cancun() Fork {
	return n.cancun
}

// RPC returns the RPC endpoint for the network.
func (n Network) RPC() (string, error) {
	switch n.id {
	case MainnetNetworkID:
		return MainnetRPC, nil
	case TestnetNetworkID:
		return TestnetRPC, nil
	case DevnetNetworkID:
		return DevnetRPC, nil
	default:
		return "", fmt.Errorf("unsupported RPC for network: %s", n.string)
	}
}

// NewNetwork returns a new Immutable zkEVM network with the specified name
// which must be one of "mainnet", "testnet", or "devnet"
func NewNetwork(name string) (*Network, error) {
	switch name {
	case "devnet":
		return &Network{
			string:      name,
			genesisJSON: immutableGenesisDevnetJSON,
			id:          DevnetNetworkID,
			cancun:      DevnetCancunFork,
		}, nil
	case "testnet":
		return &Network{
			string:      name,
			genesisJSON: immutableGenesisTestnetJSON,
			id:          TestnetNetworkID,
			cancun:      TestnetCancunFork,
		}, nil
	case "mainnet":
		return &Network{
			string:      name,
			genesisJSON: immutableGenesisMainnetJSON,
			id:          MainnetNetworkID,
			cancun:      MainnetCancunFork,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported network: %s", name)
	}
}
