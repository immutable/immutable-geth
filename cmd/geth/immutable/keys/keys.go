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

package keys

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/immutable/remote/aws"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

const (
	// EmptyPrivateKeySecretValue is the value of the private key secret before it is set
	// by the node-in-pod bootstrap operation.
	EmptyPrivateKeySecretValue = "EMPTY"
	// SecretKeyEnvVar is the environment variable used to retrieve the secret path from the environment.
	SecretKeyEnvVar = "KEY_PATH"
	// RolePlaceHolder is the placeholder for the role in the secret ID template.
	RolePlaceHolder = "{POD_ROLE}"
)

var (
	ErrPrivateKeyInvalid  = errors.New("private key is invalid")
	ErrPrivateKeyNotFound = errors.New("private key is not found")
)

// SecretIDTemplate retrieves the path to the secret key from the environment.
func SecretIDTemplate() (string, error) {
	path := os.Getenv(SecretKeyEnvVar)
	if path == "" {
		return "", fmt.Errorf("failed to retrieve secret path from environment")
	}
	return path, nil
}

// SecretID is the identifier used to store a private key (decrypted keystore) in a secure remote store.
// The secret path must contain the placeholder {POD_ROLE}, which will be populated role of the node.
// They secret path must also be scoped to the correct environment.
func SecretID(secretIDTemplate, owner string) (string, error) {
	if !strings.Contains(secretIDTemplate, RolePlaceHolder) {
		return "", fmt.Errorf("secret ID template must contain the placeholder %s", RolePlaceHolder)
	}
	return strings.Replace(secretIDTemplate, RolePlaceHolder, owner, 1), nil
}

// Store is used to retrieve and store private keys
type Store interface {
	PutSecretString(key, content string) error
	GetSecret(ctx context.Context, id string) (string, error)
}

// Render will retrieve or create a private key for a node.
// The address pertaining to the key will be logged.
func Render(
	ctx context.Context,
	store Store,
	secretID string,
) (*ecdsa.PrivateKey, error) {
	// Try and pull existing private key
	existingKeyHex, err := store.GetSecret(ctx, secretID)
	if err != nil {
		if aws.IsNotFound(err) {
			return nil, errors.Join(ErrPrivateKeyNotFound, err)
		}
		return nil, err
	}
	// Key content indicates it has not been initialized
	if existingKeyHex == EmptyPrivateKeySecretValue {
		log.Info("private key is empty, creating new key secret")
		newKey, err := crypto.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate random private key: %w", err)
		}
		newKeyHex := hex.EncodeToString(crypto.FromECDSA(newKey))
		// Push the key to the store
		if err := store.PutSecretString(secretID, newKeyHex); err != nil {
			return nil, fmt.Errorf("failed to push private key to secrets manager: %w", err)
		}
		logPublicKey(&newKey.PublicKey)
		return newKey, nil
	}
	// Decode existing private key
	log.Info("decoding existing private key")
	existingKey, err := crypto.HexToECDSA(existingKeyHex)
	if err != nil {
		return nil, errors.Join(ErrPrivateKeyInvalid, err)
	}
	logPublicKey(&existingKey.PublicKey)
	return existingKey, nil
}

func logPublicKey(pubKey *ecdsa.PublicKey) {
	log.Info("key", "public", hex.EncodeToString(crypto.FromECDSAPub(pubKey)), "address", crypto.PubkeyToAddress(*pubKey).Hex())
}

// Retrieve will retrieve a private key from the store.
func Retrieve(
	ctx context.Context,
	store Store,
	secretID string,
) (*ecdsa.PrivateKey, error) {
	keyHex, err := store.GetSecret(ctx, secretID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private key for: %w", err)
	}
	privKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		return nil, errors.Join(ErrPrivateKeyInvalid, err)
	}
	return privKey, nil
}
