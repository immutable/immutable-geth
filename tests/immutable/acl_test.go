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
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/tests/immutable/erc20"
	"github.com/ethereum/go-ethereum/tests/immutable/evm"
	"github.com/stretchr/testify/require"
)

func dynamicTx(
	t *testing.T,
	ctx context.Context,
	client *ethclient.Client,
	value *big.Int,
	keyFrom *ecdsa.PrivateKey,
	to common.Address,
) *types.Transaction {
	t.Helper()

	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)

	addrFrom := crypto.PubkeyToAddress(keyFrom.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, addrFrom)
	require.NoError(t, err)

	tx := types.MustSignNewTx(keyFrom,
		types.NewLondonSigner(chainID),
		&types.DynamicFeeTx{
			Nonce:     nonce,
			Gas:       22000,
			To:        &to,
			Value:     value,
			GasTipCap: big.NewInt(settings.PriceLimit),
			GasFeeCap: big.NewInt(settings.PriceLimit),
		})
	return tx
}

func TestImmutableACL_BlockedTransaction_ExpectedToReturnUnathorizedError(t *testing.T) {
	if config.blockedUser == nil {
		t.Skip("blockedUser/blockedprivkey was not set and is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, xerr := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, xerr)

	// Simple deploy ERC20 test
	_, _, _, err := erc20.DeployERC20(legacyTxOpts(t, ctx, client.EthClient(), config.blockedUser), client.EthClient())
	require.Error(t, err)
	require.Equal(t,
		txpool.ErrTxIsUnauthorized.Error(),
		err.Error(),
		fmt.Sprintf("unexpected error: %v, should be: %v", err, txpool.ErrTxIsUnauthorized))
}

func TestImmutableACL_Blocklist(t *testing.T) {
	if config.blockedUser == nil {
		t.Skip("blockedUser/blockedprivkey was not set and is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	// Simple deploy ERC20 test with blocked address should fail
	_, _, _, err = erc20.DeployERC20(legacyTxOpts(t, ctx, client.EthClient(), config.blockedUser), client.EthClient())
	require.Error(t, err)
	require.Equal(t,
		txpool.ErrTxIsUnauthorized.Error(),
		err.Error(),
		fmt.Sprintf("unexpected error: %v, should be: %v", err, txpool.ErrTxIsUnauthorized))

	t.Log("Deploy with blockeduser - should fail - Done")

	// Simple deploy ERC20 test with not blocked address should pass
	txOpts := legacyTxOpts(t, ctx, client.EthClient(), config.testUser)
	_, deployTx, coin, err := erc20.DeployERC20(txOpts, client.EthClient())
	require.NoError(t, err)
	deployReceipt, err := bind.WaitMined(ctx, client, deployTx)
	require.NoError(t, err)
	require.Equal(t, deployReceipt.Status, uint64(1))
	t.Logf("ERC20 Deployed to address %s with TX Hash %s", deployReceipt.ContractAddress.String(), deployTx.Hash().String())

	// Simple transfer test should also pass
	aliceKey, _ := crypto.GenerateKey()
	alice := crypto.PubkeyToAddress(aliceKey.PublicKey)

	xferTx, err := coin.Transfer(txOpts, alice, big.NewInt(1))
	require.NoError(t, err)
	_, err = bind.WaitMined(ctx, client, xferTx)
	require.NoError(t, err)
	t.Log("Transfer to alice - Done")

	// Simple native transfer from testuser to alice
	aliceTx := dynamicTx(t, ctx, client.EthClient(), big.NewInt(1), config.testUser.key, alice)
	err = client.EthClient().SendTransaction(ctx, aliceTx)
	require.NoError(t, err)
	t.Log("Transfer native to alice - should pass - Done")

	// Simple native transfer from testuser to blocked user
	tx := dynamicTx(t, ctx, client.EthClient(), big.NewInt(1), config.testUser.key, config.blockedUser.address)
	err = client.EthClient().SendTransaction(ctx, tx)
	require.Error(t, err)
	require.Equal(t,
		txpool.ErrTxIsUnauthorized.Error(),
		err.Error(),
		fmt.Sprintf("unexpected error: %v, should be: %v", err, txpool.ErrTxIsUnauthorized),
	)
	t.Log("Transfer to blockeduser - should fail - Done")
}

func TestImmutableACL_SimpleTransferTest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, xerr := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, xerr)

	// Simple transfer test should also pass
	aliceKey, _ := crypto.GenerateKey()
	alice := crypto.PubkeyToAddress(aliceKey.PublicKey)

	// Simple native transfer from testuser to alice
	aliceTx := dynamicTx(t, ctx, client.EthClient(), big.NewInt(1), config.testUser.key, alice)
	err := client.EthClient().SendTransaction(ctx, aliceTx)
	require.NoError(t, err)
	t.Log("Transfer native to alice - should pass", aliceTx.Hash())

	_, err = bind.WaitMined(ctx, client, aliceTx)
	require.NoError(t, err)
	t.Log("Transfer native to alice - should pass - Done")
}
