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
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

const scheme = rawdb.HashScheme

// All tests use full blocks as light chains are not relevant to Immutable
const full = true

// newCanonicalWithInvariants creates a chain database, and injects a deterministic canonical
// full blockchain. The chain is configured with reorgs blocked.
// The database and genesis specification for block generation
// are also returned in case more test blocks are needed later. Evolved from newCanonical
func newCanonicalWithInvariants(engine consensus.Engine, n int, full bool, scheme string) (ethdb.Database, *Genesis, *BlockChain, error) {
	chainCfg := *params.AllEthashProtocolChanges
	chainCfg.IsReorgBlocked = true
	var (
		genesis = &Genesis{
			BaseFee: big.NewInt(params.InitialBaseFee),
			Config:  &chainCfg,
		}
	)
	return newCanonicalWithGenesis(engine, n, full, scheme, genesis)
}

// testFullBlockReorgWithInvariants will generate a blockchain with no blocks. It will then create two sequences of full blocks
// of lengths defined by `first` and `second`. It will then attempt to insert the two sequences of blocks into the blockchain.
// It will then check that the chain is a validate sequence of linked hashes.
func testFullBlockReorgWithInvariants(t *testing.T, first, second []int64, expectReorg bool, expectedChainLength uint64) {
	// Create a pristine chain and database
	genDb, _, blockchain, err := newCanonicalWithInvariants(ethash.NewFaker(), 0, full, scheme)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	// Generate the first chain for inserting
	firstBlocks, _ := GenerateChain(params.TestChainConfig, blockchain.GetBlockByHash(blockchain.CurrentBlock().Hash()), ethash.NewFaker(), genDb, len(first), func(i int, b *BlockGen) {
		b.OffsetTime(first[i])
	})
	secondBlocks, _ := GenerateChain(params.TestChainConfig, blockchain.GetBlockByHash(blockchain.CurrentBlock().Hash()), ethash.NewFaker(), genDb, len(second), func(i int, b *BlockGen) {
		b.OffsetTime(second[i])
	})

	if _, err := blockchain.InsertChain(firstBlocks); err != nil {
		t.Fatalf("failed to insert firstBlocks chain: %v", err)
	}
	_, err = blockchain.InsertChain(secondBlocks)
	if expectReorg && !errors.Is(err, ErrReorgAttempted) {
		t.Fatalf("expected reorg error but got: %v", err)
	}
	if !expectReorg && err != nil {
		t.Fatalf("failed to insert secondBlocks chain: %v", err)
	}
	// Check that the chain is valid number and link wise
	prev := blockchain.CurrentBlock()
	for block := blockchain.GetBlockByNumber(blockchain.CurrentBlock().Number.Uint64() - 1); block.NumberU64() != 0; prev, block = block.Header(), blockchain.GetBlockByNumber(block.NumberU64()-1) {
		if prev.ParentHash != block.Hash() {
			t.Errorf("parent block hash mismatch: have %x, want %x", prev.ParentHash, block.Hash())
		}
	}
	chainLength := blockchain.CurrentBlock().Number.Uint64()
	if chainLength != expectedChainLength {
		t.Errorf("expected chain of length %d but got %d", expectedChainLength, chainLength)
	}
}

// TestInvariantExtendCanonicalBlocks this tests that we can extend a given chain by adding canonical blocks
// It works by generating a chain of length `i` and extending it by blocks `n`. Evolved from TestExtendCanonicalBlocksAfterMerge
// Will not trigger a reorg since the canonical chain is never reverted
func TestImmutableInvariantExtendCanonicalBlocks(t *testing.T) {
	length := 5

	// Make first chain starting from genesis
	_, _, processor, err := newCanonicalWithInvariants(ethash.NewFaker(), length, full, scheme)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	testInsertAfterMergeConfigurable(t, processor, length, 1, full, scheme, newCanonicalWithInvariants)
	testInsertAfterMergeConfigurable(t, processor, length, 10, full, scheme, newCanonicalWithInvariants)
	testInsertAfterMergeConfigurable(t, processor, length, 100, full, scheme, newCanonicalWithInvariants)
	testInsertAfterMergeConfigurable(t, processor, length, 1000, full, scheme, newCanonicalWithInvariants)
	testInsertAfterMergeConfigurable(t, processor, length, 1000, full, scheme, newCanonicalWithInvariants)
}

func TestImmutableInvariantReorg(t *testing.T) {
	// Since the blocks being inserted in the second sequence do not reorg the first sequence, we do not expect a reorg
	// and expect a final length of 6, as the second chain extends the canonical chain
	testFullBlockReorgWithInvariants(t, []int64{15, 16, 17, 18}, []int64{15, 16, 17, 18, 19, 20, 21, 22, 23, 24}, false, 10)

	// The second chain blocks attempts to reorg the last block of the first chain. It should fail and the final
	// chain should remain at the length of the initial chain
	testFullBlockReorgWithInvariants(t, []int64{15, 16, 17, 18}, []int64{15, 16, 17, -1, 18, 19, 20, 21, 22, 23}, true, 4)

	// The second chain attempts to do a deep reorg, of multiple blocks
	testFullBlockReorgWithInvariants(t, []int64{15, 16, 17, 18}, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, true, 4)

	// The second chain is different to the first chain, but since it is shorter, it will not cause a reorg
	testFullBlockReorgWithInvariants(t, []int64{15, 16, 17, 18, 19, 20}, []int64{1, 2}, false, 6)

	// Both chains are of equal length. Due to randomness, in the protocol, this will attempt to replace the head ~half the time
	// Comment it out due to test flakiness.
	//testFullBlockReorgWithInvariants(t, []int64{15, 16, 17, 18}, []int64{1, 2, 3, 4}, true, 4)

	// Short chains, where new chain extends the canonical chain
	testFullBlockReorgWithInvariants(t, []int64{15}, []int64{15, 16}, false, 2)

	// No initial chain, new chain becomes the canonical chain
	testFullBlockReorgWithInvariants(t, []int64{}, []int64{15, 16, 17}, false, 3)
}
