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

import "github.com/ethereum/go-ethereum/common"

// AddressProvider is an interface for providing a list of Ethereum addresses.
// Implementations of this interface should return a map where the keys are
// Ethereum addresses (common.Address) and the values are of an empty struct
// (struct{}{}).
//
// This interface is typically used for providing a list of
// addresses that are used for access control purposes, such as blocklists
// or allowlists.
type AddressProvider interface {
	// Provide a list of hex addressses
	Provide() map[common.Address]struct{}
}
