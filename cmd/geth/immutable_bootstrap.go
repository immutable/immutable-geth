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
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/env"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/keys"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/role"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	ethnode "github.com/ethereum/go-ethereum/node"
	"github.com/urfave/cli/v2"
)

// renderValidatorKey will create the validator's private key if it doesn't exist in the store.
// It will write the address corresponding to the key to a file at the provided dir path.
func renderValidatorKey(
	ctx context.Context,
	store keys.Store,
	dataDirpath string,
	secretID string,
) error {
	// Get or generate priv key
	privKey, err := keys.Render(
		ctx,
		store,
		secretID,
	)
	if err != nil {
		return err
	}

	// Store address on file at root data directory
	return storeAddr(privKey, filepath.Join(dataDirpath, "address"))
}

// renderP2PKey will create a new P2P key if it doesn't exist in the store.
// This is only required for boot node b/c the other nodes rely on its public key.
func renderP2PKey(
	ctx context.Context,
	store keys.Store,
	r role.Role,
	secretID string,
) error {
	if r != role.Boot {
		return fmt.Errorf("only boot node can render P2P key, got %s", r.String())
	}

	// Get or generate priv key
	_, err := keys.Render(
		ctx,
		store,
		secretID,
	)

	return err
}

// storeAddr will derive address from key and write it to a file at the provided filepath.
// This is required to have per-pod address allocation without running geth containers with a script.
func storeAddr(key *ecdsa.PrivateKey, filepath string) error {
	if key == nil {
		return fmt.Errorf("priv key is nil")
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	if err := os.WriteFile(filepath, []byte(addr.Hex()), 0644); err != nil {
		return fmt.Errorf("failed to write address to file: %w", err)
	}
	return nil
}

// renderChainState will initialize the chain and lightchain dbs with the genesis block.
// It will write files and directories to the provided root dir path.
func renderChainState(destDirpath string, genesis *core.Genesis) error {
	// Do not re-initialize, assume already initialized
	gethDirpath := filepath.Join(destDirpath, "geth")
	empty, err := isEmptyOrDoesNotExist(gethDirpath)
	if err != nil {
		return err
	}
	if !empty {
		log.Info("geth dir is not empty, skipping bootstrap", "gethDirpath", gethDirpath)
		return nil
	}
	// Instantiate node based on config and dir
	c := &cli.Context{} // NOTE: This context is empty so it can only provide default/nil values
	cfg := loadBaseConfig(c)
	cfg.Node.DataDir = destDirpath
	stack, err := ethnode.New(&cfg.Node)
	if err != nil {
		return fmt.Errorf("failed to create node: %w", err)
	}
	defer stack.Close()

	// Init genesis for chain and lightchain dbs of node
	for _, name := range []string{"chaindata", "lightchaindata"} {
		// Create and open leveldb
		chaindb, err := stack.OpenDatabaseWithFreezer(name, 0, 0, "", "", false)
		if err != nil {
			return fmt.Errorf("failed to open db: %w", err)
		}
		defer chaindb.Close()

		// Create and open trie db
		triedb := utils.MakeTrieDatabaseNoContext(cfg.Eth.StateScheme, chaindb, false, false, false)
		defer triedb.Close()

		// Write genesis block to trie and save to leveldb
		if _, _, err := core.SetupGenesisBlock(chaindb, triedb, genesis); err != nil {
			return fmt.Errorf("failed to write genesis block to db: %w", err)
		}
	}
	log.Info("rendered chain state")
	return nil
}

// ordinalFromPodName will extract the last element of the POD_NAME
// env var delimited by '-' and return it as an ordinal.
func ordinalFromPodNameEnvVar() (int, error) {
	// Read pod name for validator ordinal
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		return 0, fmt.Errorf("POD_NAME must be set")
	}
	return ordinalFromPod(podName)
}

// envFromPodNamespaceEnvVar will extract the environment from the POD_NAMESPACE
// env var.
func envFromPodNamespaceEnvVar() (env.Environment, error) {
	// Read pod namespace for the environment
	ns := os.Getenv("POD_NAMESPACE")
	if ns == "" {
		return env.Environment{}, fmt.Errorf("POD_NAMESPACE must be set")
	}
	return env.NewFromString(ns)
}

// isEmptyOrDoesNotExist returns true if the directory is empty or does not exist
func isEmptyOrDoesNotExist(dirpath string) (bool, error) {
	f, err := os.Open(dirpath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, fmt.Errorf("failed to open directory: %w", err)
	}
	defer f.Close()

	if _, err := f.Readdirnames(1); err != nil {
		if errors.Is(err, io.EOF) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
