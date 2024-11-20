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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/keys"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/node"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/role"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

const TestSecretIDTemplate = "{POD_ROLE}"

func TestImmutableBootstrap_AllRoles_CreateKeysAndChainState(t *testing.T) {
	ctx := context.Background()
	// Temp dir
	testDirpath, err := os.MkdirTemp("", "immutable-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testDirpath)

	// Configuration
	ordinal := 0

	// Subtest specification
	type tests struct {
		role     role.Role
		name     string
		store    keys.Store
		addr     string
		err      error
		secretID string
	}
	ts := []tests{}
	// These roles generate keys
	for _, r := range []role.Role{role.Validator, role.Boot} {
		// secretID :=
		ts = append(ts, []tests{
			{
				r,
				fmt.Sprintf("%sWithoutKey", r.String()),
				&mockArtefactStore{},
				"",
				keys.ErrPrivateKeyNotFound,
				"",
			},
			{
				r,
				fmt.Sprintf("%sWithKey", r.String()),
				testNewStoreWithKey(t, r, ordinal),
				"0x02F0d131F1f97aef08aEc6E3291B957d9Efe7105",
				nil,
				node.CanonicalNodeName(r, ordinal),
			},
			{
				r,
				fmt.Sprintf("%sWithEmptyKey", r.String()),
				testNewStoreWithEmptyKey(t, r, ordinal),
				"",
				keys.ErrPrivateKeyInvalid,
				node.CanonicalNodeName(r, ordinal),
			},
			{
				r,
				fmt.Sprintf("%sWithReadyKey", r.String()),
				testNewStoreWithReadyKey(t, r, ordinal),
				"",
				nil,
				node.CanonicalNodeName(r, ordinal),
			},
		}...)
	}
	// These roles don't generate keys
	for _, r := range []role.Role{role.RPC, role.Partner, role.PartnerPublic} {
		ts = append(ts, []tests{
			{
				r,
				fmt.Sprintf("%sWithoutKey", r.String()),
				&mockArtefactStore{},
				"",
				nil,
				"",
			},
		}...)
	}
	// Run subtests
	for i := range ts {
		test := ts[i]
		t.Run(test.name, func(t *testing.T) {
			// Make subdir
			subtestDirpath := filepath.Join(testDirpath, test.name)
			if err := os.Mkdir(subtestDirpath, os.ModePerm); err != nil {
				t.Fatal(err)
			}
			// Run bootstrap func
			b, err := bootstrapFactory(
				test.role,
				test.store,
				"us-east-2",
				subtestDirpath,
				test.secretID,
				core.ImmutableGenesisBlock("devnet"),
			)
			if err != nil {
				t.Fatal(err)
			}
			if err := b.Bootstrap(context.Background()); err != nil {
				if test.err == nil {
					t.Fatalf("unexpected error %v, %+v", err, test)
				}
				if !errors.Is(err, test.err) {
					t.Fatalf("expected error %v, got %v, %+v", test.err, err, test)
				}
				// Expected an error, nothing else to test
				return
			}

			// Key validation
			if test.role == role.Boot || test.role == role.Validator {
				// Now read all the secrets that should have been generated and stored and validate them
				// Key hex to compare against keystore
				secretID, err := keys.SecretID(TestSecretIDTemplate, node.CanonicalNodeName(test.role, ordinal))
				if err != nil {
					t.Fatal(err)
				}
				keyHex, err := test.store.GetSecret(ctx, secretID)
				if err != nil {
					t.Fatal(err)
				}
				// Validate secret contains valid private key
				if _, err := crypto.HexToECDSA(keyHex); err != nil {
					t.Fatal(err)
				}
			}

			// Address file validation
			if test.role == role.Validator {
				addr, err := os.ReadFile(filepath.Join(subtestDirpath, "address"))
				if err != nil {
					t.Fatal(err)
				}
				if !common.IsHexAddress(string(addr)) {
					t.Fatalf("address file does not contain a valid address %s", string(addr))
				}
				if test.addr != "" {
					if string(addr) != test.addr {
						t.Fatalf("expected address %s, got %s", test.addr, string(addr))
					}
				}
			}

			// Chain state validation
			if test.role != role.Boot {
				testRenderedChainStateDirs(t, subtestDirpath)
			}
		})
	}
}

func testNewStoreWithKey(t *testing.T, r role.Role, ordinal int) keys.Store {
	t.Helper()
	storeWithKey := &mockArtefactStore{}
	secretID, err := keys.SecretID(TestSecretIDTemplate, node.CanonicalNodeName(r, ordinal))
	if err != nil {
		t.Fatal(err)
	}
	if err := storeWithKey.PushSecretString(
		secretID,
		"48aa455c373ec5ce7fefb0e54f44a215decdc85b9047bc4d09801e038909bdbe",
	); err != nil {
		t.Fatal(err)
	}
	return storeWithKey
}

func testNewStoreWithEmptyKey(t *testing.T, r role.Role, ordinal int) keys.Store {
	t.Helper()
	storeWithNoKey := &mockArtefactStore{}
	secretID, err := keys.SecretID(TestSecretIDTemplate, node.CanonicalNodeName(r, ordinal))
	if err != nil {
		t.Fatal(err)
	}
	if err := storeWithNoKey.PushSecretString(
		secretID,
		"",
	); err != nil {
		t.Fatal(err)
	}
	return storeWithNoKey
}

func testNewStoreWithReadyKey(t *testing.T, r role.Role, ordinal int) keys.Store {
	t.Helper()
	storeWithReadyKey := &mockArtefactStore{}
	secretID, err := keys.SecretID(TestSecretIDTemplate, node.CanonicalNodeName(r, ordinal))
	if err != nil {
		t.Fatal(err)
	}
	if err := storeWithReadyKey.PushSecretString(
		secretID,
		keys.EmptyPrivateKeySecretValue,
	); err != nil {
		t.Fatal(err)
	}
	return storeWithReadyKey
}

func testRenderedChainStateDirs(t *testing.T, dirpath string) {
	// Check data dir contains expected files
	dirTests := []struct {
		dirpath string
		empty   bool
	}{
		{dirpath, false},
		{filepath.Join(dirpath, "geth"), false},
		{filepath.Join(dirpath, "geth", "chaindata"), false},
		{filepath.Join(dirpath, "geth", "lightchaindata"), false},
		{filepath.Join(dirpath, "keystore"), true},
	}
	for _, dirTest := range dirTests {
		empty, err := isEmptyOrDoesNotExist(dirTest.dirpath)
		if err != nil {
			t.Fatal(err)
		}
		if empty != dirTest.empty {
			t.Fatalf("dir %s state is incorrect", dirTest.dirpath)
		}
	}
	expectedFiles := []string{
		filepath.Join(dirpath, "geth", "LOCK"),
		filepath.Join(dirpath, "geth", "nodekey"),
	}
	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Fatalf("expected file %s to exist", expectedFile)
		}
	}
}

func TestImmutablePodEnvVars(t *testing.T) {
	tests := []struct {
		podName      string
		ordinal      int
		podNamespace string
		podValid     bool
		nsValid      bool
	}{
		{"", 0, "", false, false},
		{"zkevm-geth-validator-0", 0, "dev", true, true},
		{"zkevm-geth-validator-1", 1, "sandbox", true, true},
		{"zkevm-geth-validator-2", 2, "prod", true, true},
		{"zkevm-geth-validator-3", 3, "abc", true, false},
		{"zkevm-geth-validator-", 0, "abc", false, false},
	}

	for _, test := range tests {
		if err := os.Setenv("POD_NAME", test.podName); err != nil {
			t.Fatal(err)
		}
		if err := os.Setenv("POD_NAMESPACE", test.podNamespace); err != nil {
			t.Fatal(err)
		}
		if _, err := envFromPodNamespaceEnvVar(); err != nil && test.nsValid {
			t.Fatalf("expected pod namespace %s to be valid", test.podNamespace)
		}
		if ordinal, err := ordinalFromPodNameEnvVar(); err != nil && test.podValid {
			t.Fatalf("expected pod name %s to be valid", test.podName)
		} else {
			if ordinal != test.ordinal {
				t.Fatalf("expected ordinal %d, got %d", test.ordinal, ordinal)
			}
		}
	}
}

func TestImmutableIsEmpty(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "immutable-test")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		dirpath string
		isEmpty bool
	}{
		{
			dirpath: dir,
			isEmpty: true,
		},
		{
			dirpath: "./testdata",
			isEmpty: false,
		},
		{
			dirpath: "./does-not-exist",
			isEmpty: true,
		},
	}

	for _, test := range tests {
		t.Run(test.dirpath, func(t *testing.T) {
			isEmpty, err := isEmptyOrDoesNotExist(test.dirpath)
			if err != nil {
				t.Fatal(err)
			}
			if isEmpty != test.isEmpty {
				t.Errorf("expected %v, got %v", test.isEmpty, isEmpty)
			}
		})
	}
}
