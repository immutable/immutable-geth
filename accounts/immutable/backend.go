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

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/event"
)

// Backend implements accounts.Backend.
// It is a simple backend that only supports one account.
// The wallet is created on startup and never upated (Subscribe is noop).
type Backend struct {
	wallet accounts.Wallet
}

// NewBackend constructs a backend instance using the secret store provided.
func NewBackend(ctx context.Context, store SecretStore) (*Backend, error) {
	key, err := store.GetPrivateKey(ctx)
	if err != nil {
		return nil, err
	}
	return &Backend{
		wallet: newWallet(key),
	}, nil
}

// Wallets returns the single wallet supported by this backend.
func (b *Backend) Wallets() []accounts.Wallet {
	return []accounts.Wallet{b.wallet}
}

// Subscribe is noop for this backend because wallet never changes.
func (b *Backend) Subscribe(sink chan<- accounts.WalletEvent) event.Subscription {
	return event.NewSubscription(func(quit <-chan struct{}) error {
		<-quit
		return nil
	})
}
