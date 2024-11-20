// Copyright 2019 The go-ethereum Authors
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

package clique

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

func TestImmutableClique_ForksEnabled_BlocksValid(t *testing.T) {
	// Initialize a Clique chain with a single signer
	var (
		db     = rawdb.NewMemoryDatabase()
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr   = crypto.PubkeyToAddress(key.PublicKey)
		engine = New(params.AllCliqueProtocolChanges.Clique, db)
	)
	// Configs
	zero := uint64(0)
	shanghaiConfig := params.AllCliqueProtocolChanges
	shanghaiConfig.ShanghaiTime = &zero
	immutableConfig := params.AllCliqueProtocolChanges
	immutableConfig.ShanghaiTime = &zero
	immutableConfig.PrevrandaoTime = &zero
	immutableConfig.CancunTime = &zero

	// Test fixtures
	var tests = []struct {
		name          string
		config        *params.ChainConfig
		blockModifier func(*core.BlockGen)
	}{
		{"shanghai", immutableConfig, func(*core.BlockGen) {}},
		{"shanghaiWithdrawals", immutableConfig, func(block *core.BlockGen) {
			block.AddWithdrawal(&types.Withdrawal{
				Index:     0,
				Validator: 0,
				Address:   addr,
				Amount:    1,
			})
		}},
		{"immutable", immutableConfig, func(*core.BlockGen) {}},
		{"immutableWithdrawals", immutableConfig, func(block *core.BlockGen) {
			block.AddWithdrawal(&types.Withdrawal{
				Index:     0,
				Validator: 0,
				Address:   addr,
				Amount:    1,
			})
		}},
	}

	// Run
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			genspec := &core.Genesis{
				Config:    test.config,
				ExtraData: make([]byte, extraVanity+common.AddressLength+extraSeal),
				Alloc: map[common.Address]types.Account{
					addr: {Balance: big.NewInt(10000000000000000)},
				},
				BaseFee: big.NewInt(params.InitialBaseFee),
			}
			copy(genspec.ExtraData[extraVanity:], addr[:])

			// Generate a batch of blocks, each properly signed
			chain, _ := core.NewBlockChain(rawdb.NewMemoryDatabase(), nil, genspec, nil, engine, vm.Config{}, nil, nil)
			defer chain.Stop()

			_, blocks, _ := core.GenerateChainWithGenesis(genspec, engine, 3, func(i int, block *core.BlockGen) {
				// The chain maker doesn't have access to a chain, so the difficulty will be
				// lets unset (nil). Set it here to the correct value.
				block.SetDifficulty(diffInTurn)

				test.blockModifier(block)
			})
			for i, block := range blocks {
				header := block.Header()
				if i > 0 {
					header.ParentHash = blocks[i-1].Hash()
				}
				header.Extra = make([]byte, extraVanity+extraSeal)
				header.Difficulty = diffInTurn

				sig, _ := crypto.Sign(SealHash(header).Bytes(), key)
				copy(header.Extra[len(header.Extra)-extraSeal:], sig)
				blocks[i] = block.WithSeal(header)
			}
			// Insert the first two blocks and make sure the chain is valid
			db = rawdb.NewMemoryDatabase()
			chain, _ = core.NewBlockChain(db, nil, genspec, nil, engine, vm.Config{}, nil, nil)
			defer chain.Stop()

			if _, err := chain.InsertChain(blocks[:2]); err != nil {
				t.Fatalf("failed to insert initial blocks: %v", err)
			}
			if head := chain.CurrentBlock().Number.Uint64(); head != 2 {
				t.Fatalf("chain head mismatch: have %d, want %d", head, 2)
			}
		})
	}
}
