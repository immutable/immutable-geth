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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/cmd/immutable"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

var (
	VoteWaitTime = time.Second * 10
)

func addValidatorCommand(c *cli.Context) error {
	return processValidator(c, true)
}

func removeValidatorCommand(c *cli.Context) error {
	return processValidator(c, false)
}

func processValidator(c *cli.Context, add bool) error {
	ctx := c.Context
	voterURLs := c.StringSlice(immutable.Voters.Name)
	validatorNameStr := c.String(immutable.ValidatorAddress.Name)
	if !common.IsHexAddress(validatorNameStr) {
		return fmt.Errorf("address %s is not valid", validatorNameStr)
	}

	validator := common.HexToAddress(validatorNameStr)

	if len(voterURLs) == 0 {
		return errors.New("no voters provided")
	}

	voterClients, err := urlsToClients(ctx, voterURLs)
	if err != nil {
		return err
	}
	// Validate the operation against validator (read only)
	for i := range voterClients {
		signers, err := voterClients[i].GetSigners(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to get signers: %w", err)
		}
		// if adding a validator, the validator should not be present
		if add && slices.Contains(signers, validator) {
			return fmt.Errorf("validator %s is already part of validator set on validator %d", validator.String(), i)
		}
		// if removing a validator (i.e. !add), the validator should be present
		if !add && !slices.Contains(signers, validator) {
			return fmt.Errorf("validator %s is not part of validator set on validator %d", validator.String(), i)
		}
		log.Info("initial validator set", "validator", i, "validators", addressesToCSV(signers))
	}

	log.Info("Sleeping 15 seconds to give chance to interrupt")
	time.Sleep(time.Second * 15)

	// Execute votes against each existing validator
	for i := range voterClients {
		err := voterClients[i].Propose(ctx, validator, add)
		if err != nil {
			return fmt.Errorf("failed to propose: %w", err)
		}
		log.Info("proposing for validator", "validator", i, "add", add, "new", validator.String())
	}

	log.Info("Sleeping 10 seconds to let votes propagate and finalize")
	time.Sleep(VoteWaitTime)

	failed := false
	// Compare the new validator set to what is expected
	for i := range voterClients {
		signers, err := voterClients[i].GetSigners(ctx, nil)
		if err != nil {
			return err
		}

		// if we are adding a validator, success means that the set contains the new signer
		// if we are removing a validator, success means the set does not contain the signer
		success := (add && slices.Contains(signers, validator)) || (!add && !slices.Contains(signers, validator))
		if !success {
			failed = true
		}
		log.Info("final validator set", "validator", i, "validators", addressesToCSV(signers), "success", success)
	}
	if failed {
		return fmt.Errorf("voting failed")
	}
	return nil
}

// addressesToCSV converts a slice of common.Address to a string.
func addressesToCSV(addresses []common.Address) string {
	strAddresses := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		strAddresses = append(strAddresses, addr.String())
	}
	return strings.Join(strAddresses, ", ")
}

func urlsToClients(ctx context.Context, voterURLs []string) ([]*gethclient.Client, error) {
	voterClients := make([]*gethclient.Client, 0, len(voterURLs))
	for i := range voterURLs {
		rpcClient, err := rpc.DialContext(ctx, voterURLs[i])
		if err != nil {
			return nil, err
		}
		voterClients = append(voterClients, gethclient.New(rpcClient))
	}
	return voterClients, nil
}
