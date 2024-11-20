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
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// KeyStore implements the SecretStore interface based on a keystore file on disk.
// It is intended to be used for testing only. It replicates the existing geth keystore
// functionality in a much simpler manner because we only ever deal with a single, static keystore file.
// We use this type to verify our custom backend works as expected.
type KeyStore struct {
	key  *ecdsa.PrivateKey
	json string
	pw   string
}

// NewKeystore instantiates a new KeyStore by reading the keystore file from disk and
// decrypting it with the provided password. There must only be one keystore file in the
// provided directory.
func NewKeystore(keystoreDirpath, passwordFilepath string) (*KeyStore, error) {
	// Read keystore file
	files, err := os.ReadDir(keystoreDirpath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore dir: %w", err)
	}
	if len(files) != 1 || files[0].IsDir() {
		return nil, fmt.Errorf("keystore dir must contain exactly one file")
	}
	// Read json file
	keystoreJSON, err := os.ReadFile(filepath.Join(keystoreDirpath, files[0].Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore file: %w", err)
	}
	// Read pw file
	content, err := os.ReadFile(passwordFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read password file: %w", err)
	}
	pw := strings.TrimRight(string(content), "\r\n")

	// Decrypt keystore
	key, err := keystore.DecryptKey(keystoreJSON, pw)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt keystore: %w", err)
	}

	return &KeyStore{
		key:  key.PrivateKey,
		json: string(keystoreJSON),
		pw:   pw,
	}, nil
}

// GetPrivateKey returns the wallet's private key after retrieving its hex-encoded format from
// keystore file on disk.
func (ks *KeyStore) GetPrivateKey(ctx context.Context) (*ecdsa.PrivateKey, error) {
	return ks.key, nil
}

// JSON returns the keystore file's content as a string.
func (ks *KeyStore) JSON() string {
	return ks.json
}

// Password returns the password used to decrypt the keystore file.
func (ks *KeyStore) Password() string {
	return ks.pw
}
