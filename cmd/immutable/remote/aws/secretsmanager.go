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

package aws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/ethereum/go-ethereum/log"
)

type SecretsManager struct {
	sm *secretsmanager.SecretsManager
}

func NewSecretsManager(region string) (*SecretsManager, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("secret: failed to create session: %w", err)
	}
	man := &SecretsManager{
		sm: secretsmanager.New(sess, aws.NewConfig().WithRegion(region)),
	}
	return man, nil
}

// GetSecret returns the secret value for the given secret id
func (sm SecretsManager) GetSecret(ctx context.Context, id string) (string, error) {
	out, err := sm.sm.GetSecretValueWithContext(
		ctx,
		&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(id),
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get secret (%s): %w", id, err)
	}
	return *out.SecretString, nil
}

// GeneratePassword generates a random password
func (sm SecretsManager) GeneratePassword() (string, error) {
	// Create new
	pw, err := sm.sm.GetRandomPassword(&secretsmanager.GetRandomPasswordInput{})
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}
	return *pw.RandomPassword, nil
}

// PushSecretFile pushes a secret to the secret store
func (sm SecretsManager) PushSecretFile(secretName string, content []byte) error {
	return sm.push(&secretsmanager.CreateSecretInput{
		Name:                        aws.String(secretName),
		SecretBinary:                content,
		ForceOverwriteReplicaSecret: aws.Bool(true),
	})
}

// PutSecretString updates a secret in the secret store
func (sm SecretsManager) PutSecretString(key string, content string) error {
	if _, err := sm.sm.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(key),
		SecretString: aws.String(content),
	}); err != nil {
		return fmt.Errorf("failed to update secret (%s): %w", key, err)
	}
	return nil
}

// PushSecretString pushes a secret to the secret store
func (sm SecretsManager) PushSecretString(secretName string, content string) error {
	return sm.push(&secretsmanager.CreateSecretInput{
		Name:                        aws.String(secretName),
		SecretString:                aws.String(content),
		ForceOverwriteReplicaSecret: aws.Bool(true),
	})
}

func (sm SecretsManager) push(input *secretsmanager.CreateSecretInput) error {
	log.Info("Creating secret", "secretName", *input.Name)
	if _, err := sm.sm.CreateSecret(input); err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	return nil
}

func (sm SecretsManager) delete(secretName string) error { //nolint: unused
	if _, err := sm.sm.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(secretName),
		ForceDeleteWithoutRecovery: aws.Bool(true),
	}); err != nil && !IsNotFound(err) {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	log.Info("Waiting for deletion of secret", "secretName", secretName)
	for range time.NewTicker(time.Second * 5).C {
		if _, err := sm.sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		}); err != nil && IsNotFound(err) {
			return nil
		}
	}
	return nil
}

// IsNotFound returns true if the error is a not found error specific to AWS SDK.
func IsNotFound(err error) bool {
	var awsErr awserr.Error
	if !errors.As(err, &awsErr) {
		return false
	}
	return awsErr.Code() == secretsmanager.ErrCodeResourceNotFoundException
}
