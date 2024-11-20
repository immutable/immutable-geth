// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
)

func TestImmutableConfig_ForkChecks_Valid(t *testing.T) {
	tests := []struct {
		name                string
		chainID             *big.Int
		shanghaiTimestamp   uint64
		prevrandaoTimestamp uint64
		cancunTimestamp     uint64
	}{
		{
			name:                "mainnet",
			chainID:             big.NewInt(settings.MainnetNetworkID),
			shanghaiTimestamp:   uint64(settings.MainnetShanghaiFork.Unix()),
			prevrandaoTimestamp: uint64(settings.MainnetPrevrandaoFork.Unix()),
			cancunTimestamp:     uint64(settings.MainnetCancunFork.Unix()),
		},
		{
			name:                "testnet",
			chainID:             big.NewInt(settings.TestnetNetworkID),
			shanghaiTimestamp:   uint64(settings.TestnetShanghaiFork.Unix()),
			prevrandaoTimestamp: uint64(settings.TestnetPrevrandaoFork.Unix()),
			cancunTimestamp:     uint64(settings.TestnetCancunFork.Unix()),
		},
		{
			name:                "devnet",
			chainID:             big.NewInt(settings.DevnetNetworkID),
			shanghaiTimestamp:   uint64(settings.DevnetShanghaiFork.Unix()),
			prevrandaoTimestamp: uint64(settings.DevnetPrevrandaoFork.Unix()),
			cancunTimestamp:     uint64(settings.DevnetCancunFork.Unix()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ChainConfig{
				ChainID:        tt.chainID,
				ShanghaiTime:   &tt.shanghaiTimestamp,
				PrevrandaoTime: &tt.prevrandaoTimestamp,
				LondonBlock:    big.NewInt(0),
			}
			if !c.IsImmutableZKEVMShanghai(c.LondonBlock, tt.shanghaiTimestamp) {
				t.Errorf("expected %v to be shanghai", tt.shanghaiTimestamp)
			}
			if c.IsImmutableZKEVMShanghai(c.LondonBlock, tt.shanghaiTimestamp-10) {
				t.Errorf("expected %v to not be shanghai", tt.shanghaiTimestamp-10)
			}
			if !c.IsImmutableZKEVMPrevrandao(tt.prevrandaoTimestamp) {
				t.Errorf("expected %v to be prevrandao", tt.prevrandaoTimestamp)
			}
			if c.IsImmutableZKEVMPrevrandao(tt.prevrandaoTimestamp - 10) {
				t.Errorf("expected %v to not be prevrandao", tt.prevrandaoTimestamp-10)
			}
		})
	}
}
