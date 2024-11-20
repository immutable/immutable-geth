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

package test

import (
	"context"
	"crypto/sha256"

	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/tests/immutable/cancun/blobbasefee"
	"github.com/ethereum/go-ethereum/tests/immutable/cancun/blobhash"
	"github.com/ethereum/go-ethereum/tests/immutable/cancun/mcopy"
	selfdestructconstructor "github.com/ethereum/go-ethereum/tests/immutable/cancun/selfdestruct/constructor"
	selfdestructfunction "github.com/ethereum/go-ethereum/tests/immutable/cancun/selfdestruct/function"
	"github.com/ethereum/go-ethereum/tests/immutable/cancun/transientstorage"
	"github.com/ethereum/go-ethereum/tests/immutable/evm"
	"github.com/stretchr/testify/require"
)

func TestImmutable_Cancun_4844TransactionsDisabled(t *testing.T) {
	if skipCancun(t) {
		return
	}
	// Context and client
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	// Construct empty block txn
	emptyBlob := kzg4844.Blob{}
	emptyBlobCommit, err := kzg4844.BlobToCommitment(emptyBlob)
	require.NoError(t, err)
	emptyBlobProof, err := kzg4844.ComputeBlobProof(emptyBlob, emptyBlobCommit)
	require.NoError(t, err)
	emptyBlobVHash := kzg4844.CalcBlobHashV1(sha256.New(), &emptyBlobCommit)
	blobTx := types.BlobTx{
		Gas:        21000,
		BlobHashes: []common.Hash{emptyBlobVHash},
		Value:      uint256.NewInt(100),
		Sidecar: &types.BlobTxSidecar{
			Blobs:       []kzg4844.Blob{emptyBlob},
			Commitments: []kzg4844.Commitment{emptyBlobCommit},
			Proofs:      []kzg4844.Proof{emptyBlobProof},
		},
	}

	// Send request and expect rejection
	err = client.EthClient().SendTransaction(ctx, types.NewTx(&blobTx))
	require.Error(t, err)
	require.Equal(t, "transaction type not supported: type 3 rejected, blob transactions not supported", err.Error())
}

func TestImmutable_Cancun_4844BlobHashInstruction_ValidateOpCodeAndValue(t *testing.T) {
	if skipCancun(t) {
		return
	}
	// Context and client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	// Deploy the contract
	txOpts := EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser)
	_, deployTx, contract, err := blobhash.DeployBlobhash(txOpts, client.EthClient())
	require.NoError(t, err)
	deployReceipt, err := bind.WaitMined(ctx, client, deployTx)
	require.Equal(t, deployReceipt.Status, uint64(1))
	require.NoError(t, err)
	t.Logf("Contract deployed to address %s with TX Hash %s", deployReceipt.ContractAddress.String(), deployTx.Hash().String())

	// Transact
	opts := &bind.CallOpts{Context: ctx}
	blobHash, err := contract.BlobHash(opts, big.NewInt(1))
	require.NoError(t, err)
	require.Equal(t, [32]byte{}, blobHash)
}

func TestImmutable_Cancun_BlockHeaders(t *testing.T) {
	if skipCancun(t) {
		return
	}
	// Context and client
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	// Get latest block
	latestBlock, err := client.EthClient().BlockByNumber(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, latestBlock)

	// Check beacon root
	require.NotNil(t, latestBlock.BeaconRoot())
	require.Equal(t, common.MinHash.String(), latestBlock.BeaconRoot().String())

	// Check blob fields
	require.NotNil(t, latestBlock.BlobGasUsed())
	require.Equal(t, uint64(0), *latestBlock.BlobGasUsed())
	require.NotNil(t, latestBlock.ExcessBlobGas())
	require.Equal(t, uint64(0), *latestBlock.ExcessBlobGas())

	// Check withdrawals hash and list
	require.NotNil(t, latestBlock.Header().WithdrawalsHash)
	require.Equal(t, types.EmptyWithdrawalsHash.String(), latestBlock.Header().WithdrawalsHash.String())
	require.NotNil(t, latestBlock.Withdrawals())
	require.Empty(t, latestBlock.Withdrawals())
}

func TestImmutable_Cancun_1153TransientStorage_ValidatesTstoreAndTLoad(t *testing.T) {
	if skipCancun(t) {
		return
	}
	// Context and client
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)
	txOpts := EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser)

	// Deploy the contract
	_, deployTx, contract, err := transientstorage.DeployTransientStorage(txOpts, client.EthClient())
	require.NoError(t, err)
	deployReceipt, err := bind.WaitMined(ctx, client, deployTx)
	require.Equal(t, deployReceipt.Status, uint64(1))
	require.NoError(t, err)
	t.Logf("Contract deployed to address %s with TX Hash %s", deployReceipt.ContractAddress.String(), deployTx.Hash().String())

	// Transact
	key := big.NewInt(1)
	val := big.NewInt(2)
	tx, err := contract.TStoreLoad(txOpts, key, val)
	require.NoError(t, err)
	receipt, err := bind.WaitMined(ctx, client, tx)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))
	log, err := contract.ParseValue(*receipt.Logs[0])
	require.NoError(t, err)
	require.Equal(t, val, log.Value)
}

func TestImmutable_Cancun_7516BlobBaseFee_ValidateOpCodeAndValue(t *testing.T) {
	if skipCancun(t) {
		return
	}
	// Context and client
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)
	txOpts := EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser)

	// Deploy the contract
	_, deployTx, contract, err := blobbasefee.DeployBlobBaseFee(txOpts, client.EthClient())
	require.NoError(t, err)
	deployReceipt, err := bind.WaitMined(ctx, client, deployTx)
	require.Equal(t, deployReceipt.Status, uint64(1))
	require.NoError(t, err)
	t.Logf("Contract deployed to address %s with TX Hash %s", deployReceipt.ContractAddress.String(), deployTx.Hash().String())

	// Transact
	opts := &bind.CallOpts{Context: ctx}
	blobBaseFee, err := contract.BlobBaseFee(opts)
	require.NoError(t, err)
	require.NotNil(t, blobBaseFee)

	// Validate value as per default params
	require.Equal(t, uint64(params.BlobTxMinBlobGasprice), blobBaseFee.Uint64())
}

func TestImmutable_Cancun_5656Mcopy_ValidateOpCodeAndValue(t *testing.T) {
	if skipCancun(t) {
		return
	}
	// Context and client
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)
	txOpts := EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser)

	// Deploy the contract
	_, deployTx, contract, err := mcopy.DeployMcopy(txOpts, client.EthClient())
	require.NoError(t, err)
	deployReceipt, err := bind.WaitMined(ctx, client, deployTx)
	require.Equal(t, deployReceipt.Status, uint64(1))
	require.NoError(t, err)
	t.Logf("Contract deployed to address %s with TX Hash %s", deployReceipt.ContractAddress.String(), deployTx.Hash().String())

	// Transact
	opts := &bind.CallOpts{Context: ctx}
	mem, err := contract.MemoryCopy(opts)
	require.NoError(t, err)
	require.NotNil(t, mem)

	// Validate memory as per cancun/mcopy/mcopy.sol
	expected := [32]byte{}
	expected[31] = 0x50
	require.Equal(t, expected, mem)
}

func TestImmutable_Cancun_6780SelfDestruct_ValidateConstructorBehaviour(t *testing.T) {
	if skipCancun(t) {
		return
	}
	// Context and client
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)
	txOpts := EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser)

	// Recipient of self destruct
	recipient := genRandPubKey(t)

	// Set tx value
	txOpts.Value = big.NewInt(1)
	preBal, err := client.EthClient().BalanceAt(ctx, recipient, nil)
	require.NoError(t, err)

	// Deploy the contract
	addr, deployTx, _, err := selfdestructconstructor.DeploySelfDestructConstructor(txOpts, client.EthClient(), recipient)
	require.NoError(t, err)
	deployReceipt, err := bind.WaitMined(ctx, client, deployTx)
	require.Equal(t, deployReceipt.Status, uint64(1))
	require.NoError(t, err)
	t.Logf("Contract deployed to address %s with TX Hash %s", addr.String(), deployTx.Hash().String())

	// Check balance of recipient
	balance, err := client.EthClient().BalanceAt(ctx, recipient, nil)
	require.NoError(t, err)
	require.Equal(t, new(big.Int).Add(preBal, txOpts.Value), balance)

	// Check balance of contract
	balance, err = client.EthClient().BalanceAt(ctx, addr, nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())

	// Check code of contract is now empty
	code, err := client.EthClient().CodeAt(ctx, addr, nil)
	require.NoError(t, err)
	require.Empty(t, code)
}

func TestImmutable_Cancun_6780SelfDestruct_ValidateFunctionBehaviour(t *testing.T) {
	if skipCancun(t) {
		return
	}
	// Context and client
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)
	txOpts := EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser)

	// Recipient of self destruct
	recipient := genRandPubKey(t)

	// Deploy the contract
	addr, deployTx, contract, err := selfdestructfunction.DeploySelfDestructFunction(txOpts, client.EthClient())
	require.NoError(t, err)
	deployReceipt, err := bind.WaitMined(ctx, client, deployTx)
	require.Equal(t, deployReceipt.Status, uint64(1))
	require.NoError(t, err)

	// Fund the contract
	fundAmt := big.NewInt(1)
	tx := dynamicTx(t, ctx, client.EthClient(), fundAmt, config.testUser.key, addr)
	err = client.EthClient().SendTransaction(ctx, tx)
	require.NoError(t, err)
	receipt, err := bind.WaitMined(ctx, client, tx)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))

	// Assert balance
	balance, err := client.EthClient().BalanceAt(ctx, addr, nil)
	require.NoError(t, err)
	require.Equal(t, fundAmt, balance)

	// Self destruct
	tx, err = contract.SelfDestruct(txOpts, recipient)
	require.NoError(t, err)
	receipt, err = bind.WaitMined(ctx, client, tx)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))

	// Check balance of contract
	balance, err = client.EthClient().BalanceAt(ctx, addr, nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), balance.Int64())

	// Check balance of recipient
	balance, err = client.EthClient().BalanceAt(ctx, recipient, nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), balance)

	// Check that code is not empty
	code, err := client.EthClient().CodeAt(ctx, addr, nil)
	require.NoError(t, err)
	require.NotEmpty(t, code)
}
