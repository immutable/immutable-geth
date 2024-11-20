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
	"time"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/txpool/blobpool"
	"github.com/ethereum/go-ethereum/core/txpool/legacypool"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/urfave/cli/v2"
)

func immutableOverrides(ctx *cli.Context, cfg *gethConfig) {
	// Cannot set these flags at the same time
	utils.CheckExclusive(
		ctx,
		utils.ImmutableNetworkFlag,
		utils.OverrideShanghai,
	)
	utils.CheckExclusive(
		ctx,
		utils.ImmutableNetworkFlag,
		utils.OverridePrevrandao,
	)
	utils.CheckExclusive(
		ctx,
		utils.ImmutableNetworkFlag,
		utils.OverrideCancun,
	)
	// Set overrides based on network flag.
	if ctx.IsSet(utils.ImmutableNetworkFlag.Name) {
		genesis := core.ImmutableGenesisBlock(ctx.String(utils.ImmutableNetworkFlag.Name))
		cfg.Eth.OverrideShanghai = genesis.Config.ShanghaiTime
		cfg.Eth.OverridePrevrandao = genesis.Config.PrevrandaoTime
		cfg.Eth.OverrideCancun = genesis.Config.CancunTime
		// All overrides are handled in the genesis block, so we can terminate here
		return
	}
	// zkEVM flag was not set so lets check individual override flags for testing
	if ctx.IsSet(utils.OverridePrevrandao.Name) {
		val := ctx.Uint64(utils.OverridePrevrandao.Name)
		cfg.Eth.OverridePrevrandao = &val
	}
	if ctx.IsSet(utils.OverrideShanghai.Name) {
		val := ctx.Uint64(utils.OverrideShanghai.Name)
		cfg.Eth.OverrideShanghai = &val
	}
	if ctx.IsSet(utils.OverrideCancun.Name) {
		val := ctx.Uint64(utils.OverrideCancun.Name)
		cfg.Eth.OverrideCancun = &val
	}
}

// ImmutableEthConfig is the default content of config.toml that
// all zkEVM geth nodes should be configured to use.
func ImmutableEthConfig(chainID int) ethconfig.Config {
	return ethconfig.Config{
		SyncMode:           downloader.FullSync,
		NoPruning:          true, // See cmd/utils/flags.go:1727
		NoPrefetch:         false,
		NetworkId:          uint64(chainID),
		TxLookupLimit:      2350000,
		TransactionHistory: 2350000,
		StateScheme:        rawdb.HashScheme,
		LightPeers:         0,
		DatabaseCache:      512,
		SnapshotCache:      0,         // Not snapshot, full sync
		TrieCleanCache:     154 + 102, // See cmd/utils/flags.go:1786
		TrieDirtyCache:     256,
		TrieTimeout:        60 * time.Minute,
		FilterLogCacheSize: 32,
		Miner: miner.Config{
			GasCeil:           30000000,
			NewPayloadTimeout: 2 * time.Second,
			Recommit:          999999999 * time.Second,
			// GasPrice must be >= price limit for eth_maxPriorityFeePerGas
			GasPrice: big.NewInt(settings.PriceLimit),
		},
		TxPool: legacypool.Config{
			// PriceLimit enforces mimimum tip cap and/or tx gas price
			PriceLimit:   settings.PriceLimit,
			NoLocals:     true,
			Journal:      "transactions.rlp",
			Rejournal:    time.Hour,
			PriceBump:    10,
			AccountSlots: 16,
			GlobalSlots:  4096 + 1024, // urgent + floating queue capacity with 4:1 ratio
			AccountQueue: 64,
			GlobalQueue:  1024,
			Lifetime:     1 * time.Hour,
		},
		BlobPool: blobpool.Config{
			Datadir:   "blobpool",
			Datacap:   1, // If this is below 1, it will be set to a default of a few GB
			PriceBump: 100,
		},
		RPCGasCap:     300000000,
		RPCEVMTimeout: 5 * time.Second,
		GPO:           ethconfig.FullNodeGPO,
		RPCTxFeeCap:   0,
	}
}

func nodeConfig(enodes []*enode.Node) node.Config {
	return node.Config{
		AllowUnprotectedTxs:   true,
		InsecureUnlockAllowed: false,
		IPCPath:               "geth.ipc",
		P2P: p2p.Config{
			BootstrapNodes: enodes,
			MaxPeers:       100, // TODO: peer limit can be inferred from node count, but we need to understand surging behaviour of p2p/discovery
			NAT:            nil,
		},
	}
}
