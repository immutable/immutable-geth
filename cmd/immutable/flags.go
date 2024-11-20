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

package immutable

import (
	"github.com/urfave/cli/v2"
)

var (
	Env = &cli.StringFlag{
		Name:     "env",
		Usage:    "Network context (dev, sandbox, prod)",
		Category: ImmutableCategory,
		Required: false,
	}
	GasLimit = &cli.Uint64Flag{
		Name:     "gaslimit",
		Usage:    "Gas limit of the network",
		Value:    100000000,
		Category: ImmutableCategory,
	}
	BootCount = &cli.IntFlag{
		Name:     "boots",
		Usage:    "Number of boot nodes to bootstrap",
		Category: ImmutableCategory,
		Value:    1,
	}
	ValidatorCount = &cli.IntFlag{
		Name:     "validators",
		Usage:    "Number of validators to bootstrap",
		Category: ImmutableCategory,
		Value:    1,
	}
	RPCCount = &cli.IntFlag{
		Name:     "rpcs",
		Usage:    "Number of RPC nodes to bootstrap",
		Category: ImmutableCategory,
		Value:    1,
	}
	PartnerCount = &cli.IntFlag{
		Name:     "partners",
		Usage:    "Number of partner RPC nodes to bootstrap",
		Category: ImmutableCategory,
		Value:    1,
	}
	RPCScaleCount = &cli.IntFlag{
		Name:     "rpcs",
		Usage:    "Number of RPC nodes to scale out by",
		Category: ImmutableCategory,
		Value:    0,
	}
	PartnerScaleCount = &cli.IntFlag{
		Name:     "partners",
		Usage:    "Number of partner nodes to scale out by",
		Category: ImmutableCategory,
		Value:    0,
	}
	Region = &cli.StringFlag{
		Name:     "region",
		Usage:    "AWS region to use for remote bootstrap",
		Category: ImmutableCategory,
		Value:    "us-east-2",
	}
	BlockListFilepath = &cli.StringFlag{
		Name:     "blocklistfilepath",
		Usage:    "File path to blocklist file",
		Category: ImmutableCategory,
		Required: false,
	}
	Role = &cli.StringFlag{
		Name:     "role",
		Usage:    "Role of the node (boot, validator, rpc, partner, partner-public)",
		Category: ImmutableCategory,
		Required: true,
	}
	DataDirpath = &cli.StringFlag{
		Name:     "datadir",
		Usage:    "Data directory for node",
		Category: ImmutableCategory,
		Required: false,
	}
	ConfigFilepath = &cli.StringFlag{
		Name:     "config",
		Usage:    "Config filepath of the node",
		Category: ImmutableCategory,
		Required: true,
	}
	ManifestFilepath = &cli.StringFlag{
		Name:     "manifest",
		Usage:    "Manifest filepath of the node",
		Category: ImmutableCategory,
		Required: true,
	}
	ChainDirpath = &cli.StringFlag{
		Name:     "chain",
		Usage:    "Chain directory for node",
		Category: ImmutableCategory,
		Required: true,
	}
	Voters = &cli.StringSliceFlag{
		Name:     "voters",
		Category: ImmutableCategory,
		Usage:    "List of URLs for voting validators, e.g. http://localhost:8030/",
		Required: true,
	}
	ValidatorAddress = &cli.StringFlag{
		Name:     "validator",
		Usage:    "a validator address to vote in or out",
		Category: ImmutableCategory,
		Required: true,
	}
	PublicKey = &cli.StringFlag{
		Name:     "pubkey",
		Usage:    "public key",
		Category: ImmutableCategory,
		Required: true,
	}
)

func OverrideFlag(flag *cli.StringFlag, required bool) *cli.StringFlag {
	return &cli.StringFlag{
		Name:     flag.Name,
		Usage:    flag.Usage,
		Category: flag.Category,
		Required: required,
	}
}
