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
	"fmt"

	"github.com/ethereum/go-ethereum/cmd/immutable"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/urfave/cli/v2"
)

func runDecodeCommand(c *cli.Context) error {
	// Read encoded inputs
	pubKeyHex := c.String(immutable.PublicKey.Name)

	// Perform decode
	enodeURL := fmt.Sprintf("enode://%s@%s:%s", pubKeyHex, "127.0.0.1", "30300")
	enode, err := enode.Parse(enode.ValidSchemes, enodeURL)
	if err != nil {
		return fmt.Errorf("failed to parse enode URL: %v", err)
	}

	// Log the decoded info
	log.Info("success", "id", enode.ID())

	return nil
}
