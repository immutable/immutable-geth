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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/tests/immutable/erc721"
	"github.com/ethereum/go-ethereum/tests/immutable/evm"
	"github.com/stretchr/testify/require"
)

func TestImmutableLegacyERC721_CreateMintAndTransfer_Successful(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	testERC721(t, ctx, client.EthClient(), legacyTxOpts(t, ctx, client.EthClient(), config.testUser))
}

func TestImmutableEIP1559ERC721_CreateMintAndTransfer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)

	testERC721(t, ctx, client.EthClient(), EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser))
}

func testERC721(t *testing.T, ctx context.Context, client *ethclient.Client, txOpts *bind.TransactOpts) {
	t.Helper()

	alice := txOpts.From
	bob := common.BigToAddress(common.Big1)

	// Deploy a contract
	_, tx, nft, err := erc721.DeployERC721(txOpts, client, "bored apes", "BAYC")
	require.NoError(t, err)
	deployReceipt, err := bind.WaitMined(ctx, client, tx)
	require.Equal(t, deployReceipt.Status, uint64(1))
	require.NoError(t, err)
	t.Logf("ERC721 Deployed to address %s with TX Hash %s", deployReceipt.ContractAddress.String(), tx.Hash().String())

	// Mint tokens
	nftIDs := []*big.Int{common.Big1, common.Big2}
	for _, id := range nftIDs {
		tx, err := nft.MintNFT(txOpts, alice, id)
		require.NoError(t, err)
		receipt, err := bind.WaitMined(ctx, client, tx)
		require.NoError(t, err)
		t.Logf("Mint complete with TX Hash: %s, Block: %d", receipt.TxHash.String(), receipt.BlockNumber.Uint64())
	}

	// Transfer token
	for _, id := range nftIDs {
		tx, err := nft.TransferFrom(txOpts, alice, bob, id)
		require.NoError(t, err)
		receipt, err := bind.WaitMined(ctx, client, tx)
		require.NoError(t, err)
		t.Logf("Transfer complete with TX Hash: %s, Block: %d", receipt.TxHash.String(), receipt.BlockNumber.Uint64())
	}

	// Bad transfer
	_, err = nft.TransferFrom(txOpts, alice, bob, nftIDs[0])
	require.Error(t, err)
	require.Contains(t, err.Error(), "caller is not token owner or approved")
}
