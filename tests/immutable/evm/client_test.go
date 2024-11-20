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

package evm

import (
	"context"
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImmutableNewEVMClient_ValidArgs_Success(t *testing.T) {
	u, err := url.Parse("https://testurl")
	require.NoError(t, err)
	_, err = NewClient(context.Background(), u, 5, time.Second*2)
	require.NoError(t, err)
}

func TestImmutableNewEVMClient_ZeroRetryCount_Success(t *testing.T) {
	u, err := url.Parse("https://testurl")
	require.NoError(t, err)
	_, err = NewClient(context.Background(), u, 0, time.Second*2)
	require.NoError(t, err)
}

func TestImmutableEVMClient_FixedEstimator_ShouldReturn1GWEI(t *testing.T) {
	u, err := url.Parse("https://testurl")
	require.NoError(t, err)

	ctx := context.Background()

	client, err := NewClient(ctx, u, 5, time.Second*2)
	require.NoError(t, err)

	tip, err := client.SuggestGasTipCap(ctx)
	require.NoError(t, err)

	assert.Equal(t, big.NewInt(1000000000), tip)
}
