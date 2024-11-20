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

package gethclient

// This file introduces clique specific APIs to the gethclient. A separate client could have been created
// however, many private functions would have to be duplicated and the separation effort is not justifiable

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// GetSigners Retrieves the list of authorized signers at the specified block number.
func (ec *Client) GetSigners(ctx context.Context, blockNumber *big.Int) ([]common.Address, error) {
	var signers []common.Address
	err := ec.c.CallContext(
		ctx, &signers, "clique_getSigners", toBlockNumArg(blockNumber),
	)
	return signers, err
}

// Propose Adds a new authorization proposal that the signer will attempt to push through.
// If the auth parameter is true, the local signer votes for the given address to be included
// in the set of authorized signers.
// With auth set to false, the vote is against the address.
func (ec *Client) Propose(ctx context.Context, validator common.Address, auth bool) error {
	var result interface{}
	err := ec.c.CallContext(
		ctx, &result, "clique_propose", validator.String(), auth,
	)
	return err
}
