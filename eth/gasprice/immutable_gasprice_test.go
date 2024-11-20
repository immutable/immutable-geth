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

package gasprice

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

type testImmutableBackend struct {
	chain   *core.BlockChain
	pending bool // pending block available
	latest  int64
}

func (b *testImmutableBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	if number > testHead {
		return nil, nil
	}
	if number == rpc.EarliestBlockNumber {
		number = 0
	}
	if number == rpc.FinalizedBlockNumber {
		return b.chain.CurrentFinalBlock(), nil
	}
	if number == rpc.SafeBlockNumber {
		return b.chain.CurrentSafeBlock(), nil
	}
	if number == rpc.LatestBlockNumber {
		number = rpc.BlockNumber(b.latest)
	}
	if number == rpc.PendingBlockNumber {
		if b.pending {
			number = rpc.BlockNumber(b.latest) + 1
		} else {
			return nil, nil
		}
	}
	return b.chain.GetHeaderByNumber(uint64(number)), nil
}

func (b *testImmutableBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	if number > rpc.BlockNumber(b.latest) {
		return nil, nil
	}
	if number == rpc.EarliestBlockNumber {
		number = 0
	}
	if number == rpc.FinalizedBlockNumber {
		number = rpc.BlockNumber(b.chain.CurrentFinalBlock().Number.Uint64())
	}
	if number == rpc.SafeBlockNumber {
		number = rpc.BlockNumber(b.chain.CurrentSafeBlock().Number.Uint64())
	}
	if number == rpc.LatestBlockNumber {
		number = rpc.BlockNumber(b.latest)
	}
	if number == rpc.PendingBlockNumber {
		if b.pending {
			number = rpc.BlockNumber(b.latest) + 1
		} else {
			return nil, nil
		}
	}
	return b.chain.GetBlockByNumber(uint64(number)), nil
}

func (b *testImmutableBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return b.chain.GetReceiptsByHash(hash), nil
}

func (b *testImmutableBackend) PendingBlockAndReceipts() (*types.Block, types.Receipts) {
	if b.pending {
		block := b.chain.GetBlockByNumber(uint64(b.latest) + 1)
		return block, b.chain.GetReceiptsByHash(block.Hash())
	}
	return nil, nil
}

func (b *testImmutableBackend) ChainConfig() *params.ChainConfig {
	return b.chain.Config()
}

func (b *testImmutableBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return nil
}

func (b *testImmutableBackend) CurrentHeader() *types.Header {
	return b.chain.CurrentHeader()
}

func (b *testImmutableBackend) GetBlockByNumber(number uint64) *types.Block {
	return b.chain.GetBlockByNumber(number)
}

func (b *testImmutableBackend) teardown() {
	b.chain.Stop()
}

func newTestImmutableBackend(t *testing.T, pending bool, height int, head int, txns map[uint64]*types.DynamicFeeTx) *testImmutableBackend {
	var (
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr   = crypto.PubkeyToAddress(key.PublicKey)
		config = *params.TestChainConfig // needs copy because it is modified below
		gspec  = &core.Genesis{
			Config:  &config,
			Alloc:   types.GenesisAlloc{addr: {Balance: big.NewInt(math.MaxInt64)}},
			BaseFee: big.NewInt(49),
		}
	)
	config.ChainID = big.NewInt(settings.TestnetNetworkID)
	signer := types.LatestSigner(gspec.Config)
	londonBlock := big.NewInt(0)
	config.LondonBlock = londonBlock
	config.ArrowGlacierBlock = londonBlock
	config.GrayGlacierBlock = londonBlock
	never := uint64(100000000000000000)
	config.PrevrandaoTime = &never
	config.ShanghaiTime = nil // Cannot have shanghai in genesis (withdrawals header)
	engine := ethash.NewFaker()

	_, blocks, _ := core.GenerateChainWithGenesis(gspec, engine, height, func(i int, b *core.BlockGen) {
		b.SetCoinbase(common.Address{1})
		txdata, ok := txns[b.Number().Uint64()]
		if ok {
			b.AddTx(types.MustSignNewTx(key, signer, &types.DynamicFeeTx{
				ChainID:   gspec.Config.ChainID,
				Nonce:     b.TxNonce(addr),
				To:        &common.Address{},
				Gas:       30000,
				GasFeeCap: txdata.GasFeeCap,
				GasTipCap: txdata.GasTipCap,
				Data:      []byte{},
			}))
		}
	})

	// Construct testing chain
	chain, err := core.NewBlockChain(rawdb.NewMemoryDatabase(), &core.CacheConfig{TrieCleanNoPrefetch: true}, gspec, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create local chain, %v", err)
	}
	chain.InsertChain(blocks)
	chain.SetFinalized(chain.GetBlockByNumber(uint64(head)).Header())
	chain.SetSafe(chain.GetBlockByNumber(uint64(head)).Header())
	return &testImmutableBackend{chain: chain, pending: pending, latest: int64(head)}
}

func TestImmutablePriceLimitFloor(t *testing.T) {
	priceData := make(map[uint64]*types.DynamicFeeTx, 0)
	for i := 0; i <= 100; i++ {
		priceData[uint64(i)] = &types.DynamicFeeTx{
			GasFeeCap: big.NewInt(10 * params.GWei),
			GasTipCap: big.NewInt(10 * params.GWei),
		}
	}
	config := Config{
		Blocks:     20,
		Percentile: 60,
		Default:    big.NewInt(params.GWei),
	}
	backend := newTestImmutableBackend(t, false, 101, 100, priceData)
	oracle := NewOracle(backend, config)

	got, err := oracle.SuggestTipCap(context.Background())
	backend.teardown()
	if err != nil {
		t.Fatalf("Failed to retrieve recommended gas price: %v", err)
	}
	expected := big.NewInt(settings.PriceLimit)
	if got.Cmp(expected) != 0 {
		t.Fatalf("Gas price mismatch, want %d, got %d", expected, got)
	}
}

func TestImmutableHandleEmptyBlocks(t *testing.T) {
	priceData := make(map[uint64]*types.DynamicFeeTx, 0)
	for i := 0; i <= 50; i++ {
		priceData[uint64(i)] = &types.DynamicFeeTx{
			GasFeeCap: big.NewInt(550 * params.GWei),
			GasTipCap: big.NewInt(500 * params.GWei),
		}
	}
	config := Config{
		Blocks:     20,
		Percentile: 60,
		Default:    big.NewInt(params.GWei),
	}
	var cases = []struct {
		head   *big.Int // Head number
		expect *big.Int // Expected gasprice suggestion
	}{
		// Past 20 blocks all 500 gweis, so expect tip of 500 Gwei
		{big.NewInt(50), big.NewInt(params.GWei * int64(500))},
		// Past 20 blocks are all empty, so pull 40 blocks, 23 prices have 10 Gwei and 17 prices have 500 Gwei. 60% is top 17 prices, which is 500 Gwei.
		{big.NewInt(73), big.NewInt(params.GWei * int64(500))},
		// Past 20 blocks are all empty, so pull 40 blocks, 24 prices have 10 Gwei and 16 prices have 500 Gwei. 60% is top 17 prices, which is 10 Gwei.
		{big.NewInt(74), big.NewInt(params.GWei * int64(10))},
		{big.NewInt(75), big.NewInt(params.GWei * int64(10))},
	}
	for _, c := range cases {
		backend := newTestImmutableBackend(t, false, int(c.head.Int64())+1, int(c.head.Int64()), priceData)
		oracle := NewOracle(backend, config)

		got, err := oracle.SuggestTipCap(context.Background())
		backend.teardown()
		if err != nil {
			t.Fatalf("Failed to retrieve recommended gas price: %v", err)
		}
		if got.Cmp(c.expect) != 0 {
			t.Fatalf("Gas price mismatch, want %d, got %d", c.expect, got)
		}
	}
}
