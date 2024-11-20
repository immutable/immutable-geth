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

package accesscontrol

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestImmutableCSVProvider_load(t *testing.T) {
	tests := []struct {
		filePath      string
		expectedError bool
		// Use string here for better readability/debugging
		expectedAddress []string
	}{
		{"testdata/nonexistent.txt", true, []string{}},
		{"testdata/empty.txt", true, []string{}},
		{"testdata/gibberish.txt", true, []string{}},
		{"testdata/blocklist.txt", false, []string{
			common.HexToAddress("0x72a5843cc08275C8171E582972Aa4fDa8C397B2A").String(),
			common.HexToAddress("0x7F19720A857F834887FC9A7bC0a0fBe7Fc7f8102").String(),
			common.HexToAddress("0x1da5821544e25c636c1417ba96ade4cf6d2f9b5a").String(),
			common.HexToAddress("0x7Db418b5D567A4e0E8c59Ad71BE1FcE48f3E6107").String(),
		}},
		{"testdata/newlines.txt", false, []string{
			common.HexToAddress("0x72a5843cc08275C8171E582972Aa4fDa8C397B2A").String(),
			common.HexToAddress("0x7F19720A857F834887FC9A7bC0a0fBe7Fc7f8102").String(),
			common.HexToAddress("0x1da5821544e25c636c1417ba96ade4cf6d2f9b5a").String(),
			common.HexToAddress("0x7Db418b5D567A4e0E8c59Ad71BE1FcE48f3E6107").String(),
		}},
	}

	for _, test := range tests {
		t.Run(test.filePath, func(t *testing.T) {
			sdnProvider, err := newCSVProvider(test.filePath)

			if test.expectedError && err == nil {
				t.Errorf("Expected an error, but got none")
			} else if !test.expectedError {
				require.NoError(t, err, "Expected no error, but got an error")
				// Check if the loaded addresses match the expected addresses
				loadedAddresses := make([]string, 0, len(sdnProvider.addresses))
				for addr := range sdnProvider.addresses {
					loadedAddresses = append(loadedAddresses, addr.String())
				}
				require.ElementsMatch(t, test.expectedAddress, loadedAddresses, "Loaded addresses do not match expected addresses")
			}
		})
	}
}
