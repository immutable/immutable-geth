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
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// MockAddressProvider is a mock implementation of AddressProvider for testing.
type MockAddressProvider struct {
	addresses map[common.Address]struct{}
}

func (m *MockAddressProvider) Provide() map[common.Address]struct{} {
	return m.addresses
}

func contractCreation(nonce uint64, gaslimit uint64, key *ecdsa.PrivateKey) *types.Transaction {
	tx, _ := types.SignTx(types.NewContractCreation(nonce, big.NewInt(0), gaslimit, big.NewInt(1), nil), types.NewEIP155Signer(big.NewInt(1337)), key)
	return tx
}

func transaction(nonce uint64, gaslimit uint64, key *ecdsa.PrivateKey) *types.Transaction {
	return pricedTransaction(nonce, gaslimit, big.NewInt(1), key)
}

func pricedTransaction(nonce uint64, gaslimit uint64, gasprice *big.Int, key *ecdsa.PrivateKey) *types.Transaction {
	tx, _ := types.SignTx(types.NewTransaction(nonce, common.Address{}, big.NewInt(100), gaslimit, gasprice, nil), types.NewEIP155Signer(big.NewInt(1337)), key)
	return tx
}

func TestImmutableAccessControl_Controller_IsAllowed(t *testing.T) {
	tests := []struct {
		name            string
		providers       map[string]AddressProvider
		addressToCheck  common.Address
		isAnAllowList   bool
		expectedAllowed bool
	}{
		{
			name:          "AddressInBlockedAddresses",
			isAnAllowList: false,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: map[common.Address]struct{}{
						common.HexToAddress("0x1234567890123456789012345678901234567890"): {},
					},
				},
			},
			addressToCheck:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedAllowed: false,
		},
		{
			name:          "AddressNotInBlockedAddresses",
			isAnAllowList: false,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: map[common.Address]struct{}{},
				},
			},
			addressToCheck:  common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefab"),
			expectedAllowed: true,
		},
		{
			name:          "AddressInAllowedAddresses",
			isAnAllowList: true,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: map[common.Address]struct{}{
						common.HexToAddress("0x1234567890123456789012345678901234567890"): {},
					},
				},
			},
			addressToCheck:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expectedAllowed: true,
		},
		{
			name:          "AddressNotInAllowedAddresses",
			isAnAllowList: true,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: map[common.Address]struct{}{
						common.HexToAddress("0x1234567890123456789012345678901234567890"): {},
					},
				},
			},
			addressToCheck:  common.HexToAddress("0x11111111111111111111111111"),
			expectedAllowed: false,
		},
		{
			name:          "EmptyAddresses",
			isAnAllowList: false,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: map[common.Address]struct{}{},
				},
			},
			addressToCheck:  common.HexToAddress("0x11111111111111111111111111"),
			expectedAllowed: true,
		},
		{
			name:          "NilAddresses",
			isAnAllowList: false,
			providers: map[string]AddressProvider{
				"list": &MockAddressProvider{
					addresses: nil,
				},
			},
			addressToCheck:  common.HexToAddress("0x11111111111111111111111111"),
			expectedAllowed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := &Controller{
				providers:     test.providers,
				isAnAllowList: test.isAnAllowList,
			}

			key, _ := crypto.GenerateKey()
			tx := transaction(012, 1234, key)
			allowed := controller.IsAllowed(test.addressToCheck, tx)

			if allowed != test.expectedAllowed {
				t.Errorf("Expected allowed=%t, got allowed=%t", test.expectedAllowed, allowed)
			}
		})
	}
}
