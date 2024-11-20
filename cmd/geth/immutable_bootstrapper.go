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

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/keys"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/role"
	"github.com/ethereum/go-ethereum/core"
)

// Bootstrapper is the interface for bootstrapping a node in a k8s pod
type Bootstrapper interface {
	Bootstrap(ctx context.Context) error
}

// BootBootstrapper is the bootstrapper for boot nodes
type BootBootstrapper struct {
	store       keys.Store
	region      string
	dataDirpath string // TODO: remove field when boot migration done
	secretID    string
}

// Bootstrap bootstraps a boot node
func (bb *BootBootstrapper) Bootstrap(ctx context.Context) error {
	return renderP2PKey(
		ctx,
		bb.store,
		role.Boot,
		bb.secretID,
	)
}

// ValidatorBootstrapper is the bootstrapper for validator nodes
type ValidatorBootstrapper struct {
	store       keys.Store
	region      string
	dataDirpath string
	genesis     *core.Genesis
	secretID    string
}

// Bootstrap bootstraps a validator node
func (vb *ValidatorBootstrapper) Bootstrap(ctx context.Context) error {
	if err := renderValidatorKey(
		ctx,
		vb.store,
		vb.dataDirpath,
		vb.secretID,
	); err != nil {
		return err
	}

	return renderChainState(vb.dataDirpath, vb.genesis)
}

// RPCBootstrapper is the bootstrapper for RPC and Partner nodes
type RPCBootstrapper struct {
	dataDirpath string
	genesis     *core.Genesis
}

// Bootstrap bootstraps an RPC or Partner node
func (rb *RPCBootstrapper) Bootstrap(ctx context.Context) error {
	return renderChainState(rb.dataDirpath, rb.genesis)
}
