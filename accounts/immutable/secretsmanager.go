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
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/cmd/immutable/remote/aws"
	"github.com/ethereum/go-ethereum/crypto"
)

// AWSSecretsManagerStore implements the SecretStore interface.
// It uses an AWS Secrets Manager client to retrieve the private key intended
// to be used by a wallet.
type AWSSecretsManagerStore struct {
	sm                 *aws.SecretsManager
	validatorSecretKey string
}

// NewAWSSecretsManagerStore instantiates a new AWSSecretsManagerStore.
func NewAWSSecretsManagerStore(region, validatorSecretKey string) (*AWSSecretsManagerStore, error) {
	sm, err := aws.NewSecretsManager(region)
	if err != nil {
		return nil, err
	}
	return &AWSSecretsManagerStore{
		sm:                 sm,
		validatorSecretKey: validatorSecretKey,
	}, nil
}

// GetPrivateKey returns the wallet's private key after retrieving its hex-encoded format from
// AWS Secrets Manager.
func (s *AWSSecretsManagerStore) GetPrivateKey(ctx context.Context) (*ecdsa.PrivateKey, error) {
	// Get key from AWS Secrets Manager
	privKeyHex, err := s.sm.GetSecret(ctx, s.validatorSecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get priv key from store: %w", err)
	}
	// Decode hex key
	privKeyRaw, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode priv key hex: %w", err)
	}
	key, err := crypto.ToECDSA(privKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to convert raw priv key to ECDSA: %w", err)
	}
	return key, nil
}
