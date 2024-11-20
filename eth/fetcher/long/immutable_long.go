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

package long

import (
	"math"
	"os"

	"github.com/ethereum/go-ethereum/log"
)

const (
	longSyncEnvVar = "GETH_FLAG_IMMUTABLE_LONG_RANGE_SYNC"
)

var (
	// blockFetchDistance is increased if the relevant env var is set.
	blockFetchDistance = 256
)

func init() {
	// Enable long fetch if env var is set.
	if isLongFetch := os.Getenv(longSyncEnvVar); isLongFetch != "" {
		log.Info("Long range sync enabled")
		blockFetchDistance = math.MaxInt32
	} else {
		log.Info("Long range sync disabled")
	}
}

// BlockFetchDistance returns the maximum number of blocks that should be accepted
// when fetching blocks from peers.
// The value will be large enough for any range of blocks if the relevant env var is set.
// Otherwise it will return a reasonable value for normal operation.
func BlockFetchDistance() int {
	return blockFetchDistance
}
