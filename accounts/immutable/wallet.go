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
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// wallet implements accounts.Wallet
type wallet struct {
	key     *ecdsa.PrivateKey
	account accounts.Account
}

func newWallet(key *ecdsa.PrivateKey) *wallet {
	return &wallet{
		key: key,
		account: accounts.Account{
			Address: crypto.PubkeyToAddress(key.PublicKey),
			URL:     accounts.URL{},
		},
	}
}

// Accounts returns the single account supported by this wallet.
func (w *wallet) Accounts() []accounts.Account {
	return []accounts.Account{
		w.account,
	}
}

// Contains returns whether a particular account is or is not wrapped by this wallet instance.
func (w *wallet) Contains(account accounts.Account) bool {
	return w.account.Address == account.Address
}

// SignData signs keccak256(data). The mimetype parameter describes the type of data being signed.
func (w *wallet) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	if !w.Contains(account) {
		return nil, ErrInvalidAccount
	}
	// Sign the hash using plain ECDSA operations
	return crypto.Sign(crypto.Keccak256(data), w.key)
}

// SignText signs the hash of the given text with the given account.
func (w *wallet) SignText(account accounts.Account, text []byte) ([]byte, error) {
	if !w.Contains(account) {
		return nil, ErrInvalidAccount
	}
	// Sign the hash using plain ECDSA operations
	return crypto.Sign(accounts.TextHash(text), w.key)
}

// SignTx sends the transaction to the Immutable signer.
func (w *wallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	if !w.Contains(account) {
		return nil, ErrInvalidAccount
	}
	signer := types.LatestSignerForChainID(chainID)
	return types.SignTx(tx, signer, w.key)
}

// SignTextWithPassphrase is not implemented for Immutable signer.
func (w *wallet) SignTextWithPassphrase(account accounts.Account, passphrase string, text []byte) ([]byte, error) {
	return []byte{}, errors.New("password-operations not supported on immutable signers")
}

// SignTxWithPassphrase is not implemented for Immutable signer.
func (w *wallet) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return nil, errors.New("password-operations not supported on immutable signers")
}

// SignDataWithPassphrase is not implemented for Immutable signer.
func (w *wallet) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error) {
	return nil, errors.New("password-operations not supported on immutable signers")
}

// Open is not implemented for Immutable signer.
func (w *wallet) Open(passphrase string) error {
	return errors.New("operation not supported on immutable signers")
}

// Closed is not implemented for Immutable signer.
func (w *wallet) Close() error {
	return errors.New("operation not supported on immutable signers")
}

// Derive is not implemented for Immutable signer.
func (w *wallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	return accounts.Account{}, errors.New("operation not supported on immutable signers")
}

// SelfDerive is noop for Immutable signer.
func (w *wallet) SelfDerive(bases []accounts.DerivationPath, chain ethereum.ChainStateReader) {}

// Status is required by the accounts.Wallet interface but is not required by the Immutable signer
// because it is not used by the backend.
func (w *wallet) Status() (string, error) {
	return "unlocked", nil
}

// URL is required by the accounts.Wallet interface but is not required by the Immutable signer.
func (w *wallet) URL() accounts.URL {
	return accounts.URL{
		Scheme: "immutable",
		Path:   "",
	}
}
