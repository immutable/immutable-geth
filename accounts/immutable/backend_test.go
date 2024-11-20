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
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
)

func TestImmutableBackend_WithKeyStore_ValidAddress(t *testing.T) {
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
	// Get wallet
	wallets := b.Wallets()
	if len(wallets) != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(wallets))
	}
	if len(wallets[0].Accounts()) != 1 {
		t.Fatalf("expected 1 account, got %d", len(wallets[0].Accounts()))
	}
	// Compare wallet address
	addressBytes, err := os.ReadFile("./testdata/address.txt")
	if err != nil {
		t.Fatal(err)
	}
	if wallets[0].Accounts()[0].Address.Hex() != string(addressBytes) {
		t.Fatalf("expected address %s, got %s", string(addressBytes), wallets[0].Accounts()[0].Address.Hex())
	}
}

func TestImmutableBackend_WithKeyStore_SubscribeNoOp(t *testing.T) {
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

	// Set up subscription
	sink := make(chan accounts.WalletEvent)
	sub := b.Subscribe(sink)
	// Close subscription
	go func() {
		time.Sleep(1 * time.Second)
		sub.Unsubscribe()
	}()
	// Run subscription until closed
	for err := range sub.Err() {
		if err != nil {
			t.Fatal(err)
		}
	}
}
