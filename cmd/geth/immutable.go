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
	"github.com/ethereum/go-ethereum/cmd/immutable"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/urfave/cli/v2"
)

var (
	immutableCommand = &cli.Command{
		Name: "immutable",
		Usage: `A set of commands to manage bootstrap and deployment for Immutable infrastructure.
Local bootstrap will also run the validator node that is configured by the bootstrap.`,
		Description: "",
		Subcommands: []*cli.Command{
			{
				Name:  "bootstrap",
				Usage: "bootstrap geth network configurations",
				Subcommands: []*cli.Command{
					{
						Name:   "local",
						Usage:  "bootstrap local geth network configurations",
						Action: bootstrapLocalCommand,
						Flags: flags.Merge([]cli.Flag{
							immutable.GasLimit,
							immutable.BootCount,
							immutable.ValidatorCount,
							immutable.RPCCount,
							immutable.BlockListFilepath,
							immutable.Env,
							immutable.DataDirpath,
							utils.OverrideShanghai,
							utils.OverridePrevrandao,
							utils.OverrideCancun,
							utils.SyncModeFlag,
							configFileFlag,
							utils.GCModeFlag,
						}),
					},
					{
						Name:   "node",
						Usage:  "bootstrap node from within a k8s pod",
						Action: bootstrapK8sNodeCommand,
						Flags: flags.Merge([]cli.Flag{
							immutable.Role,
							immutable.Region,
							utils.ImmutableNetworkFlag,
							immutable.DataDirpath,
						}),
					},
					{
						Name:   "rpc",
						Usage:  "bootstrap rpc node",
						Action: bootstrapExternalNodeCommand,
						Flags: flags.Merge([]cli.Flag{
							utils.ImmutableNetworkFlag,
							immutable.DataDirpath,
						}),
					},
				},
			},
			{
				Name:  "run",
				Usage: "run a process in the geth network",
				Subcommands: []*cli.Command{
					{
						Name:   "boot",
						Usage:  "run a boot node",
						Action: runBootNodeCommand,
						Flags: flags.Merge([]cli.Flag{
							immutable.Region,
							utils.ListenPortFlag,
							utils.KeyStoreDirFlag,
						}),
					},
				},
			},
			{
				Name:   "configure",
				Usage:  "configure geth node environment before every time runtime executes",
				Action: configureNodeInPodCommand,
				Flags: flags.Merge([]cli.Flag{
					immutable.Role,
					immutable.ConfigFilepath,
					immutable.DataDirpath,
				}),
			},
			{
				Name:      "rewind",
				Usage:     "rewind chain to a specific block",
				ArgsUsage: "[? <blockHash> | <blockNum>]",
				Action:    runRewindChainCommand,
				Flags: flags.Merge([]cli.Flag{
					immutable.OverrideFlag(immutable.DataDirpath, true),
				}),
			},
			{
				Name:   "decode",
				Usage:  "decode an encoded resource",
				Action: runDecodeCommand,
				Flags: flags.Merge([]cli.Flag{
					immutable.PublicKey,
				}),
			},
			{
				Name:  "vote",
				Usage: "vote in or out new validators by connecting to multiple validators simultaneously",
				Subcommands: []*cli.Command{
					{
						Name:   "add",
						Usage:  "vote in a validator",
						Action: addValidatorCommand,
						Flags: flags.Merge([]cli.Flag{
							immutable.Voters,
							immutable.ValidatorAddress,
						}),
					},
					{
						Name:   "remove",
						Usage:  "vote out a validator",
						Action: removeValidatorCommand,
						Flags: flags.Merge([]cli.Flag{
							immutable.Voters,
							immutable.ValidatorAddress,
						}),
					},
				},
			},
		},
	}
)
