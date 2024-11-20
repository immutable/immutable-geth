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

package rewind

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestRewind_FileDoesNotExist_HistoryEmpty(t *testing.T) {
	history, err := ReadRewindHistory("test.yaml")
	require.NoError(t, err)
	require.NotNil(t, history)
	require.Empty(t, history)
}

func TestRewind_WriteRecord_HistoryNotEmptyAndRecordExists(t *testing.T) {
	// Write file
	filepath := filepath.Join(t.TempDir(), "test.yaml")
	firstHash := common.HexToHash("0x123")
	history := History{{BlockNumber: 1, BlockHash: firstHash, Timestamp: time.Now()}}
	err := WriteRewindHistory(history, filepath)
	require.NoError(t, err)

	// Write file again, with a new record
	secondHash := common.HexToHash("0x456")
	history = append(history, Record{BlockNumber: 2, BlockHash: secondHash, Timestamp: time.Now()})
	err = WriteRewindHistory(history, filepath)
	require.NoError(t, err)

	// Read file
	history, err = ReadRewindHistory(filepath)
	require.NoError(t, err)

	// Verify contents
	require.NotNil(t, history)
	require.Len(t, history, 2)
	require.True(t, history.Contains(firstHash))
	require.True(t, history.Contains(secondHash))
	require.False(t, history.Contains(common.HexToHash("0x789")))
}
