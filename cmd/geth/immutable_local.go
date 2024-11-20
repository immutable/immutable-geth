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
	"os"
	"time"

	"github.com/ethereum/go-ethereum/cmd/immutable"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

const (
	// nodePassword is used by every node to encrypt their keystore
	nodePassword = "password"
)

func bootstrapLocalCommand(c *cli.Context) error {
	// Bootstrap flags
	gasLimit := c.Uint64(immutable.GasLimit.Name)
	bootCount := c.Int(immutable.BootCount.Name)
	validatorCount := c.Int(immutable.ValidatorCount.Name)
	rpcCount := c.Int(immutable.RPCCount.Name)
	blockListFilepath := c.String(immutable.BlockListFilepath.Name)
	externalConfigFilepath := c.String(configFileFlag.Name)
	remoteNetwork := c.String(immutable.Env.Name)
	rootDirpath := c.String(immutable.DataDirpath.Name)
	if rootDirpath == "" {
		rootDirpath = os.TempDir()
	}

	// Flags fed directly into geth client
	gethFlags := []string{}
	if remoteNetwork != "" {
		gethFlags = append(gethFlags, "--zkevm", remoteNetwork)
		if validatorCount > 0 {
			return fmt.Errorf("cannot have validators in %s", remoteNetwork)
		}
	}
	if c.IsSet(utils.GCModeFlag.Name) {
		gethFlags = append(gethFlags, "--gcmode", c.String(utils.GCModeFlag.Name))
	} else {
		gethFlags = append(gethFlags, "--gcmode", "archive")
	}
	if c.IsSet(utils.OverrideShanghai.Name) {
		shanghaiTimestamp := c.Uint64(utils.OverrideShanghai.Name)
		gethFlags = append(gethFlags, "--override.shanghai", fmt.Sprint(shanghaiTimestamp))
	}
	if c.IsSet(utils.OverridePrevrandao.Name) {
		prevrandaoTimestamp := c.Uint64(utils.OverridePrevrandao.Name)
		gethFlags = append(gethFlags, "--override.prevrandao", fmt.Sprint(prevrandaoTimestamp))
	}
	if c.IsSet(utils.OverrideCancun.Name) {
		cancunTimestamp := c.Uint64(utils.OverrideCancun.Name)
		gethFlags = append(gethFlags, "--override.cancun", fmt.Sprint(cancunTimestamp))
	}
	if c.IsSet(utils.SyncModeFlag.Name) {
		syncMode := c.String(utils.SyncModeFlag.Name)
		gethFlags = append(gethFlags, "--syncmode", syncMode)
	} else {
		gethFlags = append(gethFlags, "--syncmode", "full")
	}
	if validatorCount+rpcCount > 9 {
		// The 8545/8535 ports will clash if node count is high enough.
		// This is an easy limitation to avoid, but there is no need for now.
		// Local bootstrap is only used for testing.
		return fmt.Errorf("cannot have more than 9 nodes")
	}

	// Construct the bootstrapper and run the nodes
	opts := bootstrapOptions{
		rootDirpath:             rootDirpath,
		validatorCount:          validatorCount,
		bootCount:               bootCount,
		rpcCount:                rpcCount,
		gasLimit:                gasLimit,
		blockListFilepath:       blockListFilepath,
		remoteNetwork:           remoteNetwork,
		remoteConfigFilepath:    externalConfigFilepath,
	}
	bootstrapper, err := NewLocalBootstrapper(&opts)
	if err != nil {
		return err
	}
	defer bootstrapper.Clean()

	g, _ := errgroup.WithContext(c.Context)

	boots, validators, rpcs := bootstrapper.boots, bootstrapper.validators, bootstrapper.rpcs
	for i := range boots {
		i := i // NOTE: https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			return boots[i].Run([]string{})
		})
	}
	// Wait for bootnodes
	time.Sleep(3 * time.Second) // TODO: feed an io.Writer into stdout/stderr of bootnode exes and check ready output before proceeding

	for i := range validators {
		i := i // NOTE: https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			return validators[i].Run(gethFlags)
		})
	}

	for i := range rpcs {
		i := i // NOTE: https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			return rpcs[i].Run(gethFlags)
		})
	}

	return g.Wait()
}
