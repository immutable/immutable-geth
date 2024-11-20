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

package main

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/env"
	"github.com/ethereum/go-ethereum/common"
)

func TestImmutablePremine(t *testing.T) {
	testAddr := common.HexToAddress("0x2E969d22e6654e064F461cf8B1314Cc0864a4914")
	type test struct {
		env          env.Environment
		premineCount int
		premineWei   *big.Int
	}

	tests := []test{
		{env: env.Devnet, premineCount: len(immutableDevPremineAddresses) + 1, premineWei: totalSupplyWei},
		{env: env.Testnet, premineCount: 1, premineWei: totalSupplyWei},
		{env: env.Mainnet, premineCount: 1, premineWei: totalSupplyWei},
	}

	for _, tc := range tests {
		got := immutablePremines(tc.env, testAddr)
		if len(got) != tc.premineCount {
			t.Fatalf("expected: %v, got: %v, %+v", tc.premineCount, len(got), tc)
		}
		for _, premine := range got {
			if premine.Wei.String() != tc.premineWei.String() {
				t.Fatalf("expected: %v, got: %v, %+v", tc.premineWei, premine.Wei, tc)
			}
		}
		if len(got) == 1 {
			if got[0].Address != testAddr {
				t.Fatalf("expected: %v, got: %v, %+v", testAddr, got[0].Address, tc)
			}
		}
	}
}
