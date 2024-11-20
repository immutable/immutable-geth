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
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
	"github.com/ethereum/go-ethereum/tests/immutable/erc20"
	"github.com/ethereum/go-ethereum/tests/immutable/evm"
	"github.com/stretchr/testify/require"
)

var (
	priceLimit      = big.NewInt(0).SetUint64(settings.PriceLimit)
	underPriceLimit = big.NewInt(0).SetUint64(settings.PriceLimit - 1)
)

func TestImmutablePriceLimit_SuggestTipCap_IsAbovePriceLimit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	tipCap, err := client.EthClient().SuggestGasTipCap(ctx)
	require.NoError(t, err)
	require.Equal(t, priceLimit.Uint64(), tipCap.Uint64())
	t.Log("Suggested tip cap:", tipCap)
}

func TestImmutablePriceLimit_PricedAboveLimit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	// Check that base fee is below price limit
	block, err := client.BlockNumber(ctx)
	require.NoError(t, err)
	feeHistory, err := client.EthClient().FeeHistory(ctx, 1, new(big.Int).SetUint64(block), nil)
	require.NoError(t, err)
	require.NotEmpty(t, feeHistory.BaseFee)
	require.Less(t, feeHistory.BaseFee[0].Uint64(), priceLimit.Uint64())
	t.Log("Base fee:", feeHistory.BaseFee[0].Uint64())

	// Dynamic tx above price limit but below price limit + basefee
	txOpts := EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser)
	txOpts.GasFeeCap = priceLimit
	txOpts.GasTipCap = priceLimit
	t.Log("Sending priced 1559 tx")
	testPricedTx(t, ctx, client, txOpts)

	// Legacy tx above price limit but below price limit + basefee
	txOpts = legacyTxOpts(t, ctx, client.EthClient(), config.testUser)
	txOpts.GasPrice = priceLimit
	t.Log("Sending priced legacy tx")
	testPricedTx(t, ctx, client, txOpts)
}

func TestImmutablePriceLimit_PricedBelowLimit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	// Dynamic tx below price limit
	txOpts := EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser)
	txOpts.GasFeeCap = priceLimit
	txOpts.GasTipCap = underPriceLimit
	t.Log("Sending underpriced 1559 tx")
	testUnderPricedTx(t, ctx, client, txOpts)

	// Legacy tx below price limit
	txOpts = legacyTxOpts(t, ctx, client.EthClient(), config.testUser)
	txOpts.GasPrice = underPriceLimit
	t.Log("Sending underpriced legacy tx")
	testUnderPricedTx(t, ctx, client, txOpts)
}

// testPricedTx verifies that a tx is accepted by RPC and mined
func testPricedTx(t *testing.T, ctx context.Context, client *evm.Client, txOpts *bind.TransactOpts) {
	t.Helper()

	// Send tx
	_, deployTx, _, err := erc20.DeployERC20(txOpts, client.EthClient())
	require.NoError(t, err)

	// Wait for tx to be mined
	_, err = bind.WaitMined(ctx, client, deployTx)
	require.NoError(t, err)
}

// testUnderPricedTx verifies that a tx is rejected by RPC
func testUnderPricedTx(t *testing.T, ctx context.Context, client *evm.Client, txOpts *bind.TransactOpts) {
	t.Helper()

	// Send tx
	_, _, _, err := erc20.DeployERC20(txOpts, client.EthClient())
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction underpriced")
}
