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
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/tests/immutable/erc20"
	"github.com/ethereum/go-ethereum/tests/immutable/evm"
	"github.com/stretchr/testify/require"
)

func TestImmutableLegacyERC20_CreateAndTransfer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	testERC20(t, ctx, client.EthClient(), legacyTxOpts(t, ctx, client.EthClient(), config.testUser))
}

func TestImmutableEIP1559ERC20_CreateAndTransfer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	testERC20(t, ctx, client.EthClient(), EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser))
}

func testERC20(t *testing.T, ctx context.Context, client *ethclient.Client, txOpts *bind.TransactOpts) {
	t.Helper()

	alice := txOpts.From
	bob := common.BigToAddress(common.Big1)

	// Deploy a contract
	_, deployTx, coin, err := erc20.DeployERC20(txOpts, client)
	require.NoError(t, err)
	t.Log("Deploy tx", deployTx.Hash().String())
	deployReceipt, err := bind.WaitMined(ctx, client, deployTx)
	require.Equal(t, deployReceipt.Status, uint64(1))
	require.NoError(t, err)
	t.Logf("ERC20 Deployed to address %s with TX Hash %s", deployReceipt.ContractAddress.String(), deployTx.Hash().String())

	// Check total supply non zero
	opts := &bind.CallOpts{Context: ctx}
	total, err := coin.TotalSupply(opts)
	require.NoError(t, err)
	require.Greater(t, total.Uint64(), uint64(0))

	// Get symbol
	symbol, err := coin.Symbol(opts)
	require.NoError(t, err)
	t.Logf("Coin: %s", symbol)

	// Transfer from owner to alice
	xferTx, err := coin.Transfer(txOpts, alice, total)
	require.NoError(t, err)
	receipt, err := bind.WaitMined(ctx, client, xferTx)
	require.NoError(t, err)
	t.Logf("Transfer complete with TX Hash: %s, Block: %d", receipt.TxHash.String(), receipt.BlockNumber.Uint64())

	// Check balance
	bal, err := coin.BalanceOf(opts, alice)
	require.NoError(t, err)
	require.Greater(t, bal.Uint64(), uint64(0))

	// Approve alice
	approveTx, err := coin.Approve(txOpts, alice, total)
	require.NoError(t, err)
	receipt, err = bind.WaitMined(ctx, client, approveTx)
	require.NotNil(t, receipt)
	t.Logf("Approval complete with TX Hash: %s, Block: %d", receipt.TxHash.String(), receipt.BlockNumber.Uint64())
	require.NoError(t, err)
	half := total.Div(total, common.Big2)

	// Transfer from alice to bob
	xferTx, err = coin.TransferFrom(txOpts, alice, bob, half)
	require.NoError(t, err)
	receipt, err = bind.WaitMined(ctx, client, xferTx)
	require.NoError(t, err)
	t.Logf("Transfer complete with TX Hash: %s, Block: %d", receipt.TxHash.String(), receipt.BlockNumber.Uint64())

	// Check balances after transfer
	balAlice, err := coin.BalanceOf(opts, alice)
	require.NoError(t, err)
	require.LessOrEqual(t, balAlice.Uint64(), half.Uint64())
	balBob, err := coin.BalanceOf(opts, bob)
	require.NoError(t, err)
	require.Equal(t, half.Uint64(), balBob.Uint64())
}
