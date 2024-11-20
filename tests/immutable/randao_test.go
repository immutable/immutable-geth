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
	"github.com/ethereum/go-ethereum/tests/immutable/evm"
	"github.com/ethereum/go-ethereum/tests/immutable/randao"
	"github.com/stretchr/testify/require"
)

func TestImmutable_RandaoOpCode_ShouldDeployAndReturnZero(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Client and 1559
	client, err := evm.NewClient(ctx, config.rpcURL, 1, 1)
	require.NoError(t, err)
	txOpts := EIP1559TxOpts(t, ctx, client.EthClient(), config.testUser)

	// Deploy Contract
	_, deployTx, contract, err := randao.DeployRandao(txOpts, client.EthClient())
	require.NoError(t, err)
	deployReceipt, err := bind.WaitMined(ctx, client, deployTx)
	require.NoError(t, err)
	require.Equal(t, deployReceipt.Status, uint64(1))
	t.Logf("Contract Deployed to address %s with TX Hash %s", deployReceipt.ContractAddress.String(), deployTx.Hash().String())

	// Transact
	opts := &bind.CallOpts{Context: ctx}
	rand, err := contract.Rand(opts)
	require.NoError(t, err)
	require.Equal(t, uint64(0), rand.Uint64())

	diff, err := contract.Difficulty(opts)
	require.NoError(t, err)
	require.Equal(t, uint64(0), diff.Uint64())
}
