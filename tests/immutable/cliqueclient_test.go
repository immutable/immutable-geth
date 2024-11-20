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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

func TestCliqueClient_AddSigner(t *testing.T) {
	if config.skipVoting {
		t.Skip("Skipping Test as per command-line flag")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rpcClient, err := gethrpc.DialContext(ctx, config.validatorURL.String())
	require.NoError(t, err)

	gethClient := gethclient.New(rpcClient)

	signers, err := gethClient.GetSigners(ctx, nil)
	require.NoError(t, err)

	// Should have 1 validator by default
	require.Equal(t, 1, len(signers), "expected 1 signer")

	// Propose a new validator
	newSigner := common.HexToAddress("0x7442eD1e3c9FD421F47d12A2742AfF5DaFBf43f8")
	err = gethClient.Propose(ctx, newSigner, true)
	require.NoError(t, err)

	// Wait a few seconds for a block
	log.Warn("Waiting 10s for voting to propagate")
	time.Sleep(10 * time.Second)

	signers, err = gethClient.GetSigners(ctx, nil)
	require.NoError(t, err)

	require.Equal(t, 2, len(signers), "expected 2 signers")
	if signers[0].Cmp(newSigner) != 0 && signers[1].Cmp(newSigner) != 0 {
		t.Fatalf("unexpected new signer. expected %s, got %s and %s",
			newSigner.String(),
			signers[0].String(),
			signers[1].String(),
		)
	}
}
