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
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func TestImmutableSecretsManager_SecretIsNotFound(t *testing.T) {
	tests := []struct {
		err        error
		isNotFound bool
	}{
		{
			awserr.New(secretsmanager.ErrCodeResourceNotFoundException, "Secret not found", nil),
			true,
		},
		{
			fmt.Errorf("err: %w", awserr.New(secretsmanager.ErrCodeResourceNotFoundException, "Secret not found", nil)),
			true,
		},
		{
			fmt.Errorf("bad"),
			false,
		},
	}
	for _, test := range tests {
		if IsNotFound(test.err) != test.isNotFound {
			t.Fatalf("expected %v, got %v: %s", test.isNotFound, IsNotFound(test.err), test.err.Error())
		}
	}
}
