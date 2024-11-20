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
	"path/filepath"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/rewind"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

const (
	rewindHistoryFilename = "rewind_history.yaml"
)

// rewindChainCommand is the command to rewind the chain to a specific block.
// If this is run on the validator node, it will proceed to create new blocks from
// the point of the block to which the chain is rewound.
func runRewindChainCommand(ctx *cli.Context) error {
	if ctx.Args().Len() < 1 {
		utils.Fatalf("This command requires an argument (block number or hash).")
	}

	// Set up the node handle.
	stack, _ := makeConfigNode(ctx)
	defer stack.Close()
	chain, db := utils.MakeChain(ctx, stack, false)
	defer db.Close()

	// Get the header to which the chain should be rewound.
	header, err := getRewindHeader(ctx, db)
	if err != nil {
		return err
	}

	// If the rewind has already occurred, do not repeat.
	historyFilepath := filepath.Join(stack.DataDir(), rewindHistoryFilename)
	history, err := rewind.ReadRewindHistory(historyFilepath)
	if err != nil {
		return err
	}
	log.Info("Rewind history", "history", history)
	if history.Contains(header.Hash()) {
		log.Info("Chain has already been rewound to this block (%s).", header.Hash().Hex())
		return nil
	}

	// Perform the rewind.
	log.Info("Rewinding chain to block", "block", header.Number)
	if err := chain.SetHead(header.Number.Uint64()); err != nil {
		return err
	}

	// Record the fact that the chain has been rewound.
	history = append(history, rewind.Record{
		BlockNumber: header.Number.Uint64(),
		BlockHash:   header.Hash(),
		Timestamp:   time.Now(),
	})
	return rewind.WriteRewindHistory(history, historyFilepath)
}

// getRewindHeader returns the header to which the chain should be rewound
// based on arguments passed to the command.
func getRewindHeader(ctx *cli.Context, db ethdb.Database) (*types.Header, error) {
	// Parse the argument which may be hex or a number for a block.
	arg := ctx.Args().First()
	if hashish(arg) {
		// Parse hex.
		hash := common.HexToHash(arg)
		if number := rawdb.ReadHeaderNumber(db, hash); number != nil {
			return rawdb.ReadHeader(db, hash, *number), nil
		}
		return nil, fmt.Errorf("block %x not found", hash)
	}
	// Parse number.
	number, err := strconv.ParseUint(arg, 10, 64)
	if err != nil {
		return nil, err
	}
	if hash := rawdb.ReadCanonicalHash(db, number); hash != (common.Hash{}) {
		return rawdb.ReadHeader(db, hash, number), nil
	}
	return nil, fmt.Errorf("header for block %d not found", number)
}
