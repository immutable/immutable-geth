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
	"math/big"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

// TestImmutableGenesis_ImmutableGenesis_ParsedAndValid tests that the immutable genesis blocks are valid
// so that nodes are configured correctly pre and post genesis.
func TestImmutableGenesis_ImmutableGenesis_ParsedAndValid(t *testing.T) {
	var tests = []struct {
		name       string
		addr       common.Address
		extraData  []byte
		prevrandao settings.Fork
		shanghai   settings.Fork
		cancun     settings.Fork
	}{
		{
			"mainnet",
			common.HexToAddress("dda0d9448ebe3ea43afece5fa6401f5795c19333"),
			hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000055e2ebc94b3314d387f423ef4424a95547f3f8c40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			settings.MainnetPrevrandaoFork,
			settings.MainnetShanghaiFork,
			settings.MainnetCancunFork,
		},
		{
			"testnet",
			common.HexToAddress("e567ea84e1eb3ffdc8f5aa420bf14a16eee6a809"),
			hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000512069234755cc4bb2cc59b79b7541ed5a03babb0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			settings.TestnetPrevrandaoFork,
			settings.TestnetShanghaiFork,
			settings.TestnetCancunFork,
		},
		{
			"devnet",
			common.HexToAddress("000000000013b7b1b08b3c8efe02e866f746bd38"),
			hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000092e13fb4e00a5daecf5b61553b678ad56e6985150000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			settings.DevnetPrevrandaoFork,
			settings.DevnetShanghaiFork,
			settings.DevnetCancunFork,
		},
	}
	for _, test := range tests {
		g := ImmutableGenesisBlock(test.name)
		t.Run(test.name+"genesis", func(t *testing.T) {
			assert.Equal(t, uint64(0x0), g.Nonce)
			assert.Equal(t, uint64(0x0), g.Timestamp)
			assert.Equal(t, test.extraData, g.ExtraData)
			assert.Equal(t, big.NewInt(0x1), g.Difficulty)
			assert.Equal(t, uint64(0x5f5e100), g.GasLimit)
			assert.Equal(t, common.BigToHash(common.Big0), g.Mixhash)
			assert.Equal(t, common.BigToAddress(common.Big0), g.Coinbase)
			assert.Contains(t, g.Alloc, test.addr)
			assert.Equal(t, hexutil.MustDecodeBig("0x6765c793fa10079d0000000"), g.Alloc[test.addr].Balance)
			assert.Equal(t, uint64(0), g.Number)
			assert.Equal(t, uint64(0), g.GasUsed)
			assert.Equal(t, common.BigToHash(common.Big0), g.ParentHash)
			var (
				bigNil  *big.Int
				uintNil *uint64
			)
			assert.Equal(t, bigNil, g.BaseFee)
			assert.Equal(t, uintNil, g.ExcessBlobGas)
			assert.Equal(t, uintNil, g.BlobGasUsed)

			// Make sure that genesis json matches settings.go
			assert.Equal(t, test.cancun.Unix(), int64(*g.Config.CancunTime))
			assert.Equal(t, test.prevrandao.Unix(), int64(*g.Config.PrevrandaoTime))
			assert.Equal(t, test.shanghai.Unix(), int64(*g.Config.ShanghaiTime))
		})
		t.Run(test.name+"config", func(t *testing.T) {
			c := g.Config
			assert.True(t, c.IsReorgBlocked)

			// NOTE: this test ensures that our nodes are configured to use the correct fork timetsamp!
			// Existence
			assert.NotNil(t, c.PrevrandaoTime)
			assert.NotNil(t, c.ShanghaiTime)
			// Values against human readable ones
			assert.Equal(t, test.prevrandao.Unix(), int64(*c.PrevrandaoTime))
			assert.Equal(t, test.shanghai.Unix(), int64(*c.ShanghaiTime))
			// General network config
			assert.True(t, c.IsValidImmutableZKEVM(), "Expected to be immutable zkevm")

			now := uint64(time.Now().Unix())
			// Prevrandao fork
			assert.True(t, c.IsImmutableZKEVMPrevrandao(now), test.prevrandao.IsEnabledAt(int64(now)))
			// Shanghai fork
			assert.Equal(t, c.IsShanghai(common.Big0, now), test.shanghai.IsEnabledAt(int64(now)))
			// Cancun fork
			assert.Equal(t, c.IsCancun(common.Big0, now), test.cancun.IsEnabledAt(int64(now)))
			// Check network IDs
			network, err := settings.NewNetwork(test.name)
			assert.NoError(t, err)
			assert.Equal(t, network.ID(), int(c.ChainID.Uint64()),
				"Expected network ID (%d) to match genesis chain ID (%d)", network.ID(), c.ChainID.Uint64())
		})
	}
}

func TestImmutableTestCheckCompatible(t *testing.T) {
	var tests = []struct {
		name                string
		network             string
		prevrandaoTimestamp *big.Int
		shanghaTimestamp    *big.Int
		cancunTimestamp     *big.Int
		mergeNetSplitBlock  *big.Int
	}{
		// Defaults
		{"", "devnet", nil, nil, nil, nil}, {"", "testnet", nil, nil, nil, nil}, {"", "mainnet", nil, nil, nil, nil},
		// Non-defaults
		{"devnet shanghai and cancun before merge", "devnet", big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(10000)},
		// Actual
		{"devnet actual", "devnet",
			big.NewInt(settings.DevnetPrevrandaoFork.Unix()), big.NewInt(settings.DevnetShanghaiFork.Unix()), big.NewInt(settings.DevnetCancunFork.Unix()), big.NewInt(10000)},
		{"testnet actual", "testnet",
			big.NewInt(settings.TestnetPrevrandaoFork.Unix()), big.NewInt(settings.TestnetShanghaiFork.Unix()), big.NewInt(settings.TestnetCancunFork.Unix()), big.NewInt(10000)},
		{"testnet actual", "mainnet",
			big.NewInt(settings.MainnetPrevrandaoFork.Unix()), big.NewInt(settings.MainnetShanghaiFork.Unix()), big.NewInt(settings.MainnetCancunFork.Unix()), big.NewInt(10000)},
	}

	for _, test := range tests {
		name := test.name
		if name == "" {
			name = test.network
		}
		t.Run(name, func(t *testing.T) {
			var genesis Genesis
			genesisFilepath := filepath.Join("..", "cmd", "geth", "immutable", "settings", "genesis", test.network+".json")
			err := common.LoadJSON(genesisFilepath, &genesis)
			if err != nil {
				t.Fatalf("Failed to load genesis file: %v", err)
			}
			if test.shanghaTimestamp != nil {
				ts := test.shanghaTimestamp.Uint64()
				genesis.Config.ShanghaiTime = &ts
			}
			if test.mergeNetSplitBlock != nil {
				genesis.Config.MergeNetsplitBlock = test.mergeNetSplitBlock
			}
			if genesis.Config.CheckConfigForkOrder(); err != nil {
				t.Fatalf("Failed to check config order: %v", err)
			}
		})
	}
}
