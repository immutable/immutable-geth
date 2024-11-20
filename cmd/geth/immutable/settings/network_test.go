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

package settings

import (
	"encoding/json"
	"testing"
)

func TestGenesis_GetJSON(t *testing.T) {
	var tests = []struct {
		network   string
		networkID int
	}{
		{"devnet", DevnetNetworkID},
		{"testnet", TestnetNetworkID},
		{"mainnet", MainnetNetworkID},
	}
	previousJSONs := []string{}
	for _, test := range tests {
		t.Run(test.network, func(t *testing.T) {
			network, err := NewNetwork(test.network)
			if err != nil {
				t.Fatal(err)
			}
			// ID
			if network.ID() != test.networkID {
				t.Errorf("expected network ID %d, got %d", test.networkID, network.ID())
			}
			// JSON
			if len(network.GenesisJSON()) < 1000 { // All files should be at least 1000 bytes
				t.Errorf("expected non-empty genesis JSON")
			}
			// Valid JSON
			genesis := map[string]interface{}{}
			if err := json.Unmarshal([]byte(network.GenesisJSON()), &genesis); err != nil {
				t.Fatal(err)
			}
			// Unique JSON
			if len(previousJSONs) > 0 {
				for _, previousJSON := range previousJSONs {
					if previousJSON == network.GenesisJSON() {
						t.Fatalf("expected unique genesis JSONs, got duplicate for %s and previous", test.network)
					}
				}
			}
			previousJSONs = append(previousJSONs, network.GenesisJSON())
		})
	}
	if len(previousJSONs) != len(tests) {
		t.Fatalf("expected %d unique genesis JSONs, got %d", len(tests), len(previousJSONs))
	}
}
