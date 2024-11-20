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
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/role"
	"github.com/ethereum/go-ethereum/cmd/immutable"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

var (
	// ErrOrdinalInvariant is returned when a node's ordinal is not 0 when only a single node is supported
	ErrOrdinalInvariant = fmt.Errorf("node ordinal must be 0")
	// ErrConfigRecommitInvariant is returned when config.toml has an invalid recommit value
	ErrConfigRecommitInvariant = fmt.Errorf("config.toml has invalid recommit value")

	// supportedRoles are the roles that are supported by configure command
	supportedRoles = role.AllRoles
	// singletonRoles are the roles that are only allowed to run with 1 instance
	singletonRoles = []role.Role{role.Boot, role.Validator}

	// In the case of a rotation of the single validator key, we need to remove the singleton
	// invariant for the validator role, as we need to increase the replica count
	// singletonRoles = []role.Role{role.Boot}
)

// configureNodePaths are all the filepaths that must be provided
// to the configuration of a geth node in a k8s pod
type configureNodePaths struct {
	// ConfigFilepath is the path to the node's config.toml
	ConfigFilepath string
	// DataDirpath is the path to the node's data directory
	DataDirpath string
}

func configureNodeInPodCommand(c *cli.Context) error {
	// Role
	r, err := role.NewFromString(c.String(immutable.Role.Name))
	if err != nil {
		return err
	}
	// Filepaths
	filepaths := configureNodePaths{
		ConfigFilepath: c.String(immutable.ConfigFilepath.Name),
		DataDirpath:    c.String(immutable.DataDirpath.Name),
	}

	// Env vars
	pod := os.Getenv("POD_NAME")
	if pod == "" {
		return fmt.Errorf("POD_NAME env var must be set")
	}
	return configureNodeInPod(r, pod, filepaths)
}

// configureNodeInPod will configure a geth node in a k8s (statefulset) pod based on its ordinal
func configureNodeInPod(r role.Role, pod string, filepaths configureNodePaths) error {
	// Verify that the role is supported
	if !slices.Contains(supportedRoles, r) {
		return fmt.Errorf("unsupported role: %s", r.String())
	}

	// Node ordinal
	ordinal, err := ordinalFromPod(pod)
	if err != nil {
		return err
	}
	// Verify single instances
	if slices.Contains(singletonRoles, r) {
		if ordinal != 0 {
			return ErrOrdinalInvariant
		}
	}

	// Check config.toml
	return validateConfigTOML(filepaths.ConfigFilepath)
}

// validatorConfigTOML will read the node's config.toml and check invariants
func validateConfigTOML(configFilepath string) error {
	// Open the file
	f, err := os.Open(configFilepath)
	if err != nil {
		return fmt.Errorf("failed to open toml file (%s): %w", configFilepath, err)
	}
	defer f.Close()

	// Decode the TOML
	cfg := &gethConfig{}
	if err := tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode toml file (%s): %w", configFilepath, err)
	}

	// Validate invariants
	if cfg.Eth.Miner.Recommit != ImmutableEthConfig(0).Miner.Recommit {
		return ErrConfigRecommitInvariant
	}

	// TODO: add more invariant checks here

	log.Info("validated config", "filepath", configFilepath)
	return nil
}

// ordinalFromPod extracts the ordinal (suffix) from pod name string
func ordinalFromPod(pod string) (int, error) {
	podElems := strings.Split(pod, "-")
	if len(podElems) < 2 {
		return 0, fmt.Errorf("pod name (%s) is not in the expected format", pod)
	}
	ordinal, err := strconv.Atoi(podElems[len(podElems)-1])
	if err != nil {
		return 0, fmt.Errorf("pod name (%s) suffix is not a number", pod)
	}
	return ordinal, nil
}
