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

package main

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type mockArtefactStore struct {
	secretFiles   *sync.Map
	configFiles   *sync.Map
	secretStrings *sync.Map
}

func (mas *mockArtefactStore) PushSecretFile(key string, content []byte) error {
	if mas.secretFiles == nil {
		mas.secretFiles = new(sync.Map)
	}
	mas.secretFiles.Store(key, content)
	return nil
}

func (mas *mockArtefactStore) PutSecretString(key string, content string) error {
	return mas.PushSecretString(key, content)
}

func (mas *mockArtefactStore) StoreConfigFile(filepath string, content []byte) error {
	if mas.configFiles == nil {
		mas.configFiles = new(sync.Map)
	}
	mas.configFiles.Store(filepath, content)
	return nil
}
func (mas *mockArtefactStore) PushSecretString(key, content string) error {
	if mas.secretStrings == nil {
		mas.secretStrings = new(sync.Map)
	}
	mas.secretStrings.Store(key, content)
	return nil
}
func (mas *mockArtefactStore) GetSecret(ctx context.Context, key string) (string, error) {
	if mas.secretStrings == nil {
		mas.secretStrings = new(sync.Map)
	}
	if secret, exists := mas.secretStrings.Load(key); exists {
		return secret.(string), nil
	}
	return "", awserr.New(secretsmanager.ErrCodeResourceNotFoundException, "not found", nil)
}
