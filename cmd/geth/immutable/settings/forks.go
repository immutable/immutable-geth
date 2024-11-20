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
	"time"
)

// Fork represents a human-readable date/time for a specific network fork.
// It can be used to validate unix timestamps (e.g. in genesis.json) so that
// it is easier to validate the correctness of the fork timestamps.
type Fork struct {
	time.Time
}

func newFork(timestamp string) Fork {
	const (
		forkTimestampLocation = "UTC"
		forkTimestampLayout   = time.UnixDate
	)
	loc, err := time.LoadLocation(forkTimestampLocation)
	if err != nil {
		panic(err)
	}
	t, err := time.ParseInLocation(forkTimestampLayout, timestamp, loc)
	if err != nil {
		panic(err)
	}

	if t.UTC().Format(forkTimestampLayout) != timestamp {
		panic(fmt.Sprintf("fork time %s does not match supplied %s", t.UTC().Format(forkTimestampLayout), timestamp))
	}

	return Fork{t}
}

// IsEnabledAt returns true if the given unix timestamp is past the fork timestamp
func (fork Fork) IsEnabledAt(time int64) bool {
	return time >= fork.Unix()
}
