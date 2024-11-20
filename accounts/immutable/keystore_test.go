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
	"encoding/json"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestImmutableKeyStore_ReadKeyStoreFile_ValidKeystore(t *testing.T) {
	ks, err := NewKeystore("./testdata/keystore", "./testdata/password.txt")
	if err != nil {
		t.Fatal(err)
	}
	key, err := ks.GetPrivateKey(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Check address
	addr := crypto.PubkeyToAddress(key.PublicKey)
	addressBytes, err := os.ReadFile("./testdata/address.txt")
	if err != nil {
		t.Fatal(err)
	}
	if addr.Hex() != string(addressBytes) {
		t.Fatalf("expected address %s, got %s", string(addressBytes), addr)
	}

	// Check JSON
	if err := json.Unmarshal([]byte(ks.JSON()), &struct{}{}); err != nil {
		t.Fatal(err)
	}

	// Check pw
	pwBytes, err := os.ReadFile("./testdata/password.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(pwBytes) != ks.Password() {
		t.Fatalf("expected pw %s, got %s", string(pwBytes), ks.Password())
	}
}
