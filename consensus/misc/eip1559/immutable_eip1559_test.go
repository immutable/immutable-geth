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

package eip1559

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func TestImmutableCalcBaseFee(t *testing.T) {
	tests := []struct {
		parentBaseFee   int64
		parentGasLimit  uint64
		parentGasUsed   uint64
		expectedBaseFee int64
	}{
		{params.InitialBaseFee, 20000000, 10000000, params.InitialBaseFee}, // usage == target
		{params.InitialBaseFee, 20000000, 9000000, 998000000},              // usage below target
		{params.InitialBaseFee, 20000000, 11000000, 1002000000},            // usage above target
	}
	for i, test := range tests {
		parent := &types.Header{
			Number:   common.Big32,
			GasLimit: test.parentGasLimit,
			GasUsed:  test.parentGasUsed,
			BaseFee:  big.NewInt(test.parentBaseFee),
		}
		// Set config that modifies base fee calculation
		// (cannot import core pkg for genesis due to import cycle)
		c := config()
		c.Clique = &params.CliqueConfig{
			Period: settings.SecondsPerBlock,
		}
		c.IsReorgBlocked = true
		c.ChainID = big.NewInt(settings.DevnetNetworkID)
		if !c.IsValidImmutableZKEVM() {
			t.Fatalf("test %d: invalid immutable zkevm config", i)
		}
		if have, want := CalcBaseFee(c, parent), big.NewInt(test.expectedBaseFee); have.Cmp(want) != 0 {
			t.Fatalf("test %d: have %d  want %d, ", i, have, want)
		}
	}
}
