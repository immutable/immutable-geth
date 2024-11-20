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

package immutable

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-cmp/cmp"
)

func TestImmutableWallet_FromKeystore_SignaturesNoErrors(t *testing.T) {
	// Initialize keystore
	ks, err := NewKeystore("./testdata/keystore", "./testdata/password.txt")
	if err != nil {
		t.Fatal(err)
	}
	// Initialize backend
	b, err := NewBackend(context.Background(), ks)
	if err != nil {
		t.Fatal(err)
	}
	// Get wallet and sign data
	w := b.Wallets()[0]

	tests := []struct {
		data    []byte
		tx      *types.Transaction
		chainID *big.Int
	}{
		{
			data:    []byte{1, 2, 3},
			tx:      types.NewTx(&types.LegacyTx{}),
			chainID: big.NewInt(1),
		},
	}
	for _, test := range tests {
		if _, err := w.SignData(w.Accounts()[0], "", test.data); err != nil {
			t.Fatal(err)
		}
		if _, err := w.SignText(w.Accounts()[0], test.data); err != nil {
			t.Fatal(err)
		}
		if _, err := w.SignTx(w.Accounts()[0], test.tx, test.chainID); err != nil {
			t.Fatal(err)
		}
	}
}

func TestImmutableWallet_FromKeystore_SignaturesValid(t *testing.T) {
	// Get default keystore wallet
	ks := keystore.NewKeyStore("./testdata/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	pwBytes, err := os.ReadFile("./testdata/password.txt")
	if err != nil {
		t.Fatal(err)
	}
	pw := strings.TrimRight(string(pwBytes), "\r\n")
	if err := ks.Unlock(ks.Accounts()[0], pw); err != nil {
		t.Fatal(err)
	}
	if len(ks.Accounts()) != 1 {
		t.Fatalf("expected 1 account, got %d", len(ks.Accounts()))
	}
	if len(ks.Wallets()) != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(ks.Wallets()))
	}

	// Get immutable wallet
	immutableKeystore, err := NewKeystore("./testdata/keystore", "./testdata/password.txt")
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewBackend(context.Background(), immutableKeystore)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		walletLeft  accounts.Wallet
		walletRight accounts.Wallet
	}{
		{
			walletLeft:  ks.Wallets()[0],
			walletRight: b.Wallets()[0],
		},
	}
	for _, test := range tests {
		// Sign data
		data := []byte{1, 2, 3}
		signatureLeft, err := test.walletLeft.SignData(test.walletLeft.Accounts()[0], "", data)
		if err != nil {
			t.Fatal(err)
		}
		signatureRight, err := test.walletRight.SignData(test.walletRight.Accounts()[0], "", data)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(signatureLeft, signatureRight) {
			t.Fatalf("expected signature %s, got %s", string(signatureLeft), string(signatureRight))
		}
		// Sign text
		signatureLeft, err = test.walletLeft.SignText(test.walletLeft.Accounts()[0], data)
		if err != nil {
			t.Fatal(err)
		}
		signatureRight, err = test.walletRight.SignText(test.walletRight.Accounts()[0], data)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(signatureLeft, signatureRight) {
			t.Fatalf("expected signature %s, got %s", string(signatureLeft), string(signatureRight))
		}
		// Sign tx
		tx := types.NewTx(&types.LegacyTx{})
		chainID := big.NewInt(1)
		txLeft, err := test.walletLeft.SignTx(test.walletLeft.Accounts()[0], tx, chainID)
		if err != nil {
			t.Fatal(err)
		}
		txRight, err := test.walletRight.SignTx(test.walletRight.Accounts()[0], tx, chainID)
		if err != nil {
			t.Fatal(err)
		}
		lv, lr, ls := txLeft.RawSignatureValues()
		rv, rr, rs := txRight.RawSignatureValues()
		if !cmp.Equal(lv.String(), rv.String()) {
			t.Fatalf("expected v %d, got %d", lv, rv)
		}
		if !cmp.Equal(lr.String(), rr.String()) {
			t.Fatalf("expected r %d, got %d", lr, rr)
		}
		if !cmp.Equal(ls.String(), rs.String()) {
			t.Fatalf("expected s %d, got %d", ls, rs)
		}
	}
}
