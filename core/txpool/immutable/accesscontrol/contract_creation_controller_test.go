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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestImmutableAccessControl_ContractCreationController_IsAllowed(t *testing.T) {
	key, _ := crypto.GenerateKey()
	tests := []struct {
		name            string
		providers       map[string]AddressProvider
		addressToCheck  common.Address
		tx              *types.Transaction
		isAnAllowList   bool
		expectedAllowed bool
	}{
		{
			name:          "Address is in block list and tx is of contract creation",
			isAnAllowList: false,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: map[common.Address]struct{}{
						common.HexToAddress("0x1234567890123456789012345678901234567890"): {},
					},
				},
			},
			tx:              contractCreation(1234, 123, key),
			addressToCheck:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedAllowed: false,
		},
		{
			name:          "Address is in block list and tx is not of contract creation",
			isAnAllowList: false,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: map[common.Address]struct{}{
						common.HexToAddress("0x1234567890123456789012345678901234567890"): {},
					},
				},
			},
			tx:              transaction(1234, 123, key),
			addressToCheck:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedAllowed: true,
		},
		{
			name:          "Empty blocklist but tx is of contract creation",
			isAnAllowList: false,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: map[common.Address]struct{}{},
				},
			},
			tx:              contractCreation(1234, 123, key),
			addressToCheck:  common.HexToAddress("0x11111111111111111111111111"),
			expectedAllowed: true,
		},
		{
			name:          "Nil blocklist sets but tx is of contract creation",
			isAnAllowList: false,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: nil,
				},
			},
			tx:              contractCreation(1234, 123, key),
			addressToCheck:  common.HexToAddress("0x11111111111111111111111111"),
			expectedAllowed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := &ContractCreation{
				providers:     test.providers,
				isAnAllowList: test.isAnAllowList,
			}

			allowed := controller.IsAllowed(test.addressToCheck, test.tx)

			if allowed != test.expectedAllowed {
				t.Errorf("Expected allowed=%t, got allowed=%t", test.expectedAllowed, allowed)
			}
		})
	}
}
