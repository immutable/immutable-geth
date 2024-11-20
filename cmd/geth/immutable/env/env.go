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

package env

import (
	"fmt"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
)

// Environment represents an Immutable zkEVM environment
// and assists referencing resources correctly, such as AWS secrets, k8s YAMLs, and
// configuration files.
type Environment struct {
	string
}

func (e Environment) String() string {
	return e.string
}

// ChainID returns the chain ID for the environment
func (e Environment) ChainID() int {
	switch e {
	case Devnet:
		return settings.DevnetNetworkID
	case Testnet:
		return settings.TestnetNetworkID
	case Mainnet:
		return settings.MainnetNetworkID
	default:
		panic(fmt.Sprintf("unsupported env: %s", e))
	}
}

var (
	Devnet  = Environment{"dev"}
	Testnet = Environment{"sandbox"}
	Mainnet = Environment{"prod"}
)

// NewFromString returns an Environment from a string
func NewFromString(env string) (Environment, error) {
	switch env {
	case Devnet.String():
		return Devnet, nil
	case Testnet.String():
		return Testnet, nil
	case Mainnet.String():
		return Mainnet, nil
	default:
		return Environment{}, fmt.Errorf("unsupported env: %s", env)
	}
}
