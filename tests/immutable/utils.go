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
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

// Wallet represents an account/wallet which holds currency
type Wallet struct {
	key     *ecdsa.PrivateKey
	address common.Address
}

func loadWalletFromFile(filepath string) (*Wallet, error) {
	privKeyFileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("priv key read file: %w", err)
	}
	privKeyFileBytes = bytes.ReplaceAll(privKeyFileBytes, []byte{'\n'}, []byte{})
	return loadWallet(string(privKeyFileBytes))
}

func loadWallet(privKeyHex string) (*Wallet, error) {
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("priv key from bytes: %w", err)
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)
	return &Wallet{key: privKey, address: address}, nil
}

func legacyTxOpts(t *testing.T, ctx context.Context, client *ethclient.Client, user *Wallet) *bind.TransactOpts {
	t.Helper()

	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)

	opts, err := bind.NewKeyedTransactorWithChainID(user.key, chainID)
	require.NoError(t, err)

	// Set gasprice above price limit
	opts.GasPrice = big.NewInt(0).SetUint64(settings.PriceLimit * 2)
	return opts
}

func EIP1559TxOpts(t *testing.T, ctx context.Context, client *ethclient.Client, user *Wallet) *bind.TransactOpts {
	t.Helper()

	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)

	return &bind.TransactOpts{
		Signer: func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
			return types.SignTx(t, types.NewLondonSigner(chainID), user.key)
		},
		From: user.address,
		// Nonce nil, otherwise we need to increment or call PendingNonce
		// between transactions ourselves
		Nonce: nil,
		// Leaving these as nil allows go-ethereum to calculate the values
		// based on suggested prices from the chain
		GasTipCap: nil,
		GasFeeCap: nil,
	}
}

func genRandPubKey(t *testing.T) common.Address {
	t.Helper()
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
