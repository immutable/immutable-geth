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
	"fmt"
	"testing"
	"time"
)

func TestForks_IsPastTimestamp(t *testing.T) {
	// Test certain timestamps
	var tests = []struct {
		unixTimestamp int64
		isPast        bool
		forks         []Fork
	}{
		{0, false, []Fork{
			DevnetShanghaiFork, TestnetShanghaiFork, MainnetShanghaiFork,
			DevnetPrevrandaoFork, TestnetPrevrandaoFork, MainnetPrevrandaoFork,
			DevnetCancunFork, TestnetCancunFork, MainnetCancunFork,
		}},
		// Enabled now
		{time.Now().Unix(), true, []Fork{
			DevnetShanghaiFork, TestnetShanghaiFork, MainnetShanghaiFork,
			DevnetPrevrandaoFork, TestnetPrevrandaoFork, MainnetPrevrandaoFork,
			DevnetCancunFork, TestnetCancunFork, MainnetCancunFork,
		}},
		// Scheduled forks that are not yet enabled (may be empty list)
		{time.Now().Unix(), false, []Fork{}},
	}
	for _, test := range tests {
		for _, fork := range test.forks {
			name := fmt.Sprintf("unix %d %s", test.unixTimestamp, fork)
			t.Run(name, func(t *testing.T) {
				if got, want := fork.IsEnabledAt(test.unixTimestamp), test.isPast; got != want {
					t.Errorf("got %v, want %v", got, want)
				}
				if !fork.IsEnabledAt(fork.Unix()) {
					t.Errorf("fork timestamp should be past itself")
				}
				if fork.IsEnabledAt(fork.Unix() - 10) {
					t.Errorf("fork timestamp should not be past itself - 10")
				}
			})
		}
	}
}
