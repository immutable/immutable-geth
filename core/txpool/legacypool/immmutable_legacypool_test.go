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

package legacypool

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
)

type MockAccessController struct {
	isBlocklist bool
	isAllowed   bool
}

func (m *MockAccessController) IsBlocklist() bool {
	return m.isBlocklist
}

func (m *MockAccessController) IsAllowed(from common.Address, tx *types.Transaction) bool {
	return m.isAllowed
}

func contractCreation(nonce uint64, gaslimit uint64, key *ecdsa.PrivateKey) *types.Transaction {
	tx, _ := types.SignTx(types.NewContractCreation(nonce, big.NewInt(0), gaslimit, big.NewInt(1), nil), types.HomesteadSigner{}, key)
	return tx
}

func pricedTransactionTo(nonce uint64, gaslimit uint64, gasprice *big.Int, key *ecdsa.PrivateKey, to common.Address) *types.Transaction {
	tx, _ := types.SignTx(types.NewTransaction(nonce, to, big.NewInt(100), gaslimit, gasprice, nil), types.HomesteadSigner{}, key)
	return tx
}

func TestImmutableFilterWithError(t *testing.T) {
	signer := types.NewEIP155Signer(big.NewInt(1337))
	testCases := []struct {
		description string
		pool        *LegacyPool
		expectedErr error
	}{
		{
			description: "Blocked by Blocklist",
			pool: &LegacyPool{
				accessControllers: []txpool.AccessController{
					&MockAccessController{isBlocklist: true, isAllowed: false},
					&MockAccessController{isBlocklist: true, isAllowed: false},
				},
				signer: signer,
			},
			expectedErr: txpool.ErrTxIsUnauthorized,
		},
		{
			description: "No access controllers, so it should pass",
			pool: &LegacyPool{
				accessControllers: []txpool.AccessController{},
				signer:            signer,
			},
			expectedErr: nil,
		},
		{
			description: "An access controller with only a single block list that allows",
			pool: &LegacyPool{
				accessControllers: []txpool.AccessController{
					&MockAccessController{isBlocklist: true, isAllowed: true},
				},
				signer: signer,
			},
			expectedErr: nil,
		},
	}

	key, _ := crypto.GenerateKey()
	to := common.HexToAddress("0x1111")

	tx, err := types.SignTx(types.NewTransaction(uint64(10), to, big.NewInt(1000), params.TxGas, big.NewInt(params.InitialBaseFee), nil), signer, key)
	if err != nil {
		t.Errorf("Couldn't create test transaction: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.pool.FilterWithError(tx)
			if err != tc.expectedErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

type filterWithErrorTest struct {
	name          string
	expectedError error
	inputTx       *types.Transaction
}

func TestImmutableFilterWithErrorBlocklist(t *testing.T) {
	// Fixed keys
	sdnKey, _ := crypto.GenerateKey()
	sdnAddress := crypto.PubkeyToAddress(sdnKey.PublicKey)

	randomKey, _ := crypto.GenerateKey()

	tests := []filterWithErrorTest{

		{
			name:          "BlockedTxBecauseSenderAreInBlocklist",
			expectedError: txpool.ErrTxIsUnauthorized,
			inputTx:       pricedDataTransaction(1, 100000, big.NewInt(1), sdnKey, 50),
		},
		{
			name:          "BlockedTxBecauseReceiverInBlocklist",
			expectedError: txpool.ErrTxIsUnauthorized,
			inputTx:       pricedTransactionTo(1, 100000, big.NewInt(1), randomKey, sdnAddress),
		},
		{
			name:          "NormalTransaction",
			expectedError: nil,
			inputTx:       pricedDataTransaction(1, 100000, big.NewInt(1), randomKey, 50),
		},
	}
	// Setup test
	statedb, _ := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := newTestBlockChain(params.TestChainConfig, 1000000, statedb, new(event.Feed))
	blocklist := filepath.Join(t.TempDir(), "blocklist.txt")
	blocklistFile, err := os.Create(blocklist)
	if err != nil {
		t.Fatal("Failed to create temporary blocklist")
	}

	_, err = blocklistFile.WriteString(crypto.PubkeyToAddress(sdnKey.PublicKey).Hex())
	if err != nil {
		t.Fatal("Failed to write to blocklist file")
	}
	config := testTxPoolConfig
	config.BlockListFilePaths = []string{blocklist}

	pool := New(config, blockchain)
	pool.Init(config.PriceLimit, blockchain.CurrentBlock(), makeAddressReserver())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := pool.FilterWithError(tc.inputTx); err != tc.expectedError {
				t.Errorf("Test %s: expected error %v, got %v", tc.name, tc.expectedError, err)
			}
		})
	}
}

func TestImmutableACLFileWithNewLinesShouldBeBlocked(t *testing.T) {
	statedb, _ := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := newTestBlockChain(params.TestChainConfig, 1000000, statedb, new(event.Feed))

	key, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()
	key3, _ := crypto.GenerateKey()

	keys := []*ecdsa.PrivateKey{key, key2, key3}

	addresses := []any{
		crypto.PubkeyToAddress(key.PublicKey).Hex(),
		crypto.PubkeyToAddress(key2.PublicKey).Hex(),
		crypto.PubkeyToAddress(key3.PublicKey).Hex(),
	}

	blocklist := filepath.Join(t.TempDir(), "blocklist.txt")
	file, err := os.Create(blocklist)
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}

	_, err = fmt.Fprintf(file, "%s,\n%s,\n%s", addresses...)
	if err != nil {
		t.Fatal("Failed to write to file")
	}

	config := testTxPoolConfig
	config.BlockListFilePaths = []string{blocklist}

	pool := New(config, blockchain)
	pool.Init(config.PriceLimit, blockchain.CurrentBlock(), makeAddressReserver())

	for _, key := range keys {
		tx := contractCreation(1, 100000, key)
		if err := pool.FilterWithError(tx); err != txpool.ErrTxIsUnauthorized {
			t.Fatalf("Transaction was not blocked: %v", err)
		}
	}
}

func TestImmutableFilterWithErrorBlockedTransaction(t *testing.T) {
	statedb, _ := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := newTestBlockChain(params.TestChainConfig, 1000000, statedb, new(event.Feed))

	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey)

	blocklist := filepath.Join(t.TempDir(), "blocklist.txt")
	file, err := os.Create(blocklist)
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}

	_, err = file.WriteString(address.Hex())
	if err != nil {
		t.Fatal("Failed to write to file")
	}

	config := testTxPoolConfig
	config.BlockListFilePaths = []string{blocklist}

	pool := New(config, blockchain)
	pool.Init(config.PriceLimit, blockchain.CurrentBlock(), makeAddressReserver())

	tx := contractCreation(1, 100000, key)

	if err := pool.FilterWithError(tx); err != txpool.ErrTxIsUnauthorized {
		t.Fatalf("Transaction was not blocked: %v", err)
	}
}

func TestImmutableMinGasPriceEnforced(t *testing.T) {
	t.Parallel()

	// Create the pool to test the pricing enforcement with
	statedb, _ := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := newTestBlockChain(eip1559Config, 10000000, statedb, new(event.Feed))

	txPoolConfig := DefaultConfig
	txPoolConfig.NoLocals = true
	pool := New(txPoolConfig, blockchain)
	pool.Init(txPoolConfig.PriceLimit, blockchain.CurrentBlock(), makeAddressReserver())
	defer pool.Close()

	key, _ := crypto.GenerateKey()
	testAddBalance(pool, crypto.PubkeyToAddress(key.PublicKey), big.NewInt(1000000))

	tx := pricedTransaction(0, 100000, big.NewInt(2), key)
	pool.SetGasTip(big.NewInt(tx.GasPrice().Int64() + 1))

	if err := pool.addLocal(tx); !errors.Is(err, txpool.ErrUnderpriced) {
		t.Fatalf("Min tip not enforced")
	}

	if err := pool.Add([]*types.Transaction{tx}, true, false)[0]; !errors.Is(err, txpool.ErrUnderpriced) {
		t.Fatalf("Min tip not enforced")
	}

	tx = dynamicFeeTx(0, 100000, big.NewInt(3), big.NewInt(2), key)
	pool.SetGasTip(big.NewInt(tx.GasTipCap().Int64() + 1))

	if err := pool.addLocal(tx); !errors.Is(err, txpool.ErrUnderpriced) {
		t.Fatalf("Min tip not enforced")
	}

	if err := pool.Add([]*types.Transaction{tx}, true, false)[0]; !errors.Is(err, txpool.ErrUnderpriced) {
		t.Fatalf("Min tip not enforced")
	}
	// Make sure the tx is accepted if locals are enabled
	pool.config.NoLocals = false
	if err := pool.Add([]*types.Transaction{tx}, true, false)[0]; err != nil {
		t.Fatalf("Min tip enforced with locals enabled, error: %v", err)
	}
}

func TestImmutableRebroadcasting(t *testing.T) {
	t.Parallel()

	// Create a test account and fund it
	pool, key := setupPool()
	defer pool.Close()

	account := crypto.PubkeyToAddress(key.PublicKey)
	testAddBalance(pool, account, big.NewInt(1000000))

	// Keep track of transaction events to ensure all executables get announced
	events := make(chan core.NewTxsEvent, testTxPoolConfig.AccountQueue+5)
	sub := pool.txFeed.Subscribe(events)
	defer sub.Unsubscribe()

	// Create 2 pending transactions and 1 queued transaction. pending txs should be broadcast
	pool.addRemotesSync([]*types.Transaction{
		transaction(0, 100000, key),
		transaction(1, 100000, key),
		transaction(4, 100000, key),
	})

	pending, queued := pool.Stats()
	if pending != 2 {
		t.Fatalf("pending transactions mismatched: have %d, want %d", pending, 2)
	}
	if queued != 1 {
		t.Fatalf("queued transactions mismatched: have %d, want %d", queued, 1)
	}
	if err := validateEvents(events, 2); err != nil {
		t.Fatalf("add tx event firing failed: %v", err)
	}
	if err := validatePoolInternals(pool); err != nil {
		t.Fatalf("pool internal state corrupted: %v", err)
	}

	// Reset is called by txpool.go when a new block is received.
	// We simulate 'a new block' by calling it here
	pool.Reset(nil, nil)

	pending, queued = pool.Stats()
	if pending != 2 {
		t.Fatalf("pending transactions mismatched: have %d, want %d", pending, 2)
	}
	if queued != 1 {
		t.Fatalf("queued transactions mismatched: have %d, want %d", queued, 1)
	}
	// The pending transactions should be rebroadcast. Queued transactions should not.
	if err := validateEvents(events, 2); err != nil {
		t.Fatalf("rebroadcast event firing failed: %v", err)
	}
	if err := validatePoolInternals(pool); err != nil {
		t.Fatalf("pool internal state corrupted: %v", err)
	}

	// Adding a new pending transaction should only result in the new transaction being broadcast,
	// since there was no new block
	pool.addRemotesSync([]*types.Transaction{
		transaction(2, 100000, key),
	})
	pending, queued = pool.Stats()
	if pending != 3 {
		t.Fatalf("pending transactions mismatched: have %d, want %d", pending, 3)
	}
	if queued != 1 {
		t.Fatalf("queued transactions mismatched: have %d, want %d", queued, 1)
	}
	// Only the new transaction should be broadcast
	if err := validateEvents(events, 1); err != nil {
		t.Fatalf("add tx event firing failed: %v", err)
	}
	if err := validatePoolInternals(pool); err != nil {
		t.Fatalf("pool internal state corrupted: %v", err)
	}

	// A new block should result in rebroadcast of all 3 pending transactions
	pool.Reset(nil, nil)

	pending, queued = pool.Stats()
	if pending != 3 {
		t.Fatalf("pending transactions mismatched: have %d, want %d", pending, 3)
	}
	if queued != 1 {
		t.Fatalf("queued transactions mismatched: have %d, want %d", queued, 1)
	}
	// The pending transactions should be rebroadcast. Queued transactions should not.
	if err := validateEvents(events, 3); err != nil {
		t.Fatalf("rebroadcast event firing failed: %v", err)
	}
	if err := validatePoolInternals(pool); err != nil {
		t.Fatalf("pool internal state corrupted: %v", err)
	}
}
