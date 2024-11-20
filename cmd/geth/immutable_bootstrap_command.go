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

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/keys"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/node"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/role"
	"github.com/ethereum/go-ethereum/cmd/immutable"
	"github.com/ethereum/go-ethereum/cmd/immutable/remote/aws"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/urfave/cli/v2"
)

// getSecretID will return the secret ID for the node.
// The secret ID will be based on the role and ordinal of the node.
func getSecretID(role role.Role) (string, error) {
	// Ordinal.
	ordinal, err := ordinalFromPodNameEnvVar()
	if err != nil {
		return "", err
	}
	// Secret ID template.
	secretIDTemplate, err := keys.SecretIDTemplate()
	if err != nil {
		return "", err
	}
	// Secret ID.
	secretID, err := keys.SecretID(secretIDTemplate, node.CanonicalNodeName(role, ordinal))
	if err != nil {
		return "", err
	}
	return secretID, nil
}

func bootstrapK8sNodeCommand(c *cli.Context) error {
	// Role
	r, err := role.NewFromString(c.String(immutable.Role.Name))
	if err != nil {
		return err
	}

	// Setup store
	store, err := aws.NewSecretsManager(c.String(immutable.Region.Name))
	if err != nil {
		return err
	}

	// Only retrieve secret IDs for boot and validator nodes. Other nodes do not store secret material.
	var secretID string
	switch r {
	case role.Boot, role.Validator:
		secretID, err = getSecretID(r)
		if err != nil {
			return err
		}
	default:
		secretID = ""
	}

	// Get bootstrapper and run it
	b, err := bootstrapFactory(
		r,
		store,
		c.String(immutable.Region.Name),
		c.String(immutable.DataDirpath.Name),
		secretID,
		core.ImmutableGenesisBlock(c.String(utils.ImmutableNetworkFlag.Name)),
	)
	if err != nil {
		return err
	}
	return b.Bootstrap(c.Context)
}

func bootstrapExternalNodeCommand(c *cli.Context) error {
	// Get bootstrapper and run it
	b, err := bootstrapFactory(
		role.RPC,
		nil,
		"",
		c.String(immutable.DataDirpath.Name),
		"",
		core.ImmutableGenesisBlock(c.String(utils.ImmutableNetworkFlag.Name)),
	)
	if err != nil {
		return err
	}
	return b.Bootstrap(c.Context)
}

// bootstrapFactory will create a bootstrapper based on the role
// of the relevant node. Each node has different bootstrap requirements.
func bootstrapFactory(
	r role.Role,
	store keys.Store,
	region string,
	dataDirpath string,
	secretID string,
	genesis *core.Genesis,
) (Bootstrapper, error) {
	switch r {
	case role.Boot:
		return &BootBootstrapper{
			store:       store,
			region:      region,
			dataDirpath: dataDirpath,
			secretID:    secretID,
		}, nil
	case role.Validator:
		return &ValidatorBootstrapper{
			store:       store,
			region:      region,
			dataDirpath: dataDirpath,
			genesis:     genesis,
			secretID:    secretID,
		}, nil
	case role.Partner, role.RPC, role.PartnerPublic:
		return &RPCBootstrapper{
			dataDirpath: dataDirpath,
			genesis:     genesis,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported role %s", r.String())
	}
}
