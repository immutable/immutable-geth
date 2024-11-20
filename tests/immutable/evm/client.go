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

package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"time"

	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sethvargo/go-retry"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ethClient  *ethclient.Client
	rpcClient  *gethrpc.Client
	retryCount uint64
	retryWait  time.Duration
	// Capable of giving EIP-1559 gas tip estimations according to some internal
	// strategy. Receives the ethclient.Client to fetch gas information.
	gasTipEstimator GasTipEstimator
}

const gethExecutionReverted = "execution reverted"

// TODO: Remove edgeExecutionReverted once edge matches up error strings with geth
const edgeExecutionReverted = "execution was reverted"

// TODO: Figure out a better way to handle error from hardhat
const hardhatExecutionReverted = "Error: Transaction reverted without a reason string"

var ErrExecutionReverted = errors.New("execution reverted")

// NewClient create default instance of an EVM Client
func NewClient(ctx context.Context, u *url.URL, retryCount uint64, retryWait time.Duration) (*Client, error) {
	evmClient, err := ethclient.DialContext(ctx, u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to dial context for eth client: %w", err)
	}
	rpcClient, err := gethrpc.DialContext(ctx, u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to dial context for rpc client: %w", err)
	}

	return &Client{
		ethClient:  evmClient,
		rpcClient:  rpcClient,
		retryCount: retryCount,
		retryWait:  retryWait,
		// Suggested minimum per EIP-1559
		gasTipEstimator: NewFixedGasTipEstimator(big.NewInt(1000000000)),
	}, nil
}

// FilterLogs Ethereum RPC FilterLogs
func (d *Client) FilterLogs(ctx context.Context, filter ethereum.FilterQuery) ([]types.Log, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.FilterLogs")
	var logs []types.Log
	err := d.retryDo(ctx, func(ctx context.Context) error {
		var err error
		logs, err = d.ethClient.FilterLogs(ctx, filter)
		if err != nil {
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to filter eth logs: %w", err)
	}
	return logs, nil
}

// BlockNumber Ethereum RPC BlockNumber
func (d *Client) BlockNumber(ctx context.Context) (uint64, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.BlockNumber")
	var blockNum uint64
	err := d.retryDo(ctx, func(ctx context.Context) error {
		var err error
		blockNum, err = d.ethClient.BlockNumber(ctx)
		if err != nil {
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}
	return blockNum, nil
}

// BlockByNumber Ethereum RPC BlockByNumber
func (d *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.BlockByNumber")
	var block *types.Block
	err := d.retryDo(ctx, func(ctx context.Context) error {
		var err error
		block, err = d.ethClient.BlockByNumber(ctx, number)
		if err != nil {
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by block number %s: %w", number.String(), err)
	}
	return block, nil
}

// TransactionReceipt Ethereum RPC TransactionReceipt
func (d *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.TransactionReceipt")
	var receipt *types.Receipt
	err := d.retryDo(ctx, func(ctx context.Context) error {
		var err error
		receipt, err = d.ethClient.TransactionReceipt(ctx, txHash)
		if err != nil {
			// Don't log "not found" errors, they are expected
			if err.Error() == "not found" {
				return retry.RetryableError(err)
			}
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction recept for transaction %s: %w", txHash.Hex(), err)
	}
	return receipt, nil
}

// CodeAt calls the Ethereum RPC CodeAt function.
func (d *Client) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.CodeAt")
	var codeAt []byte
	err := d.retryDo(ctx, func(ctx context.Context) error {
		var err error
		codeAt, err = d.ethClient.CodeAt(ctx, contract, blockNumber)
		if err != nil {
			// NOTE: We have to do string comparison because generated bindings don't return a type we can easily work with
			// This error is not retryable.
			if err.Error() == gethExecutionReverted || err.Error() == edgeExecutionReverted || err.Error() == hardhatExecutionReverted {
				return ErrExecutionReverted
			}
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get code at for address %s: %w", contract.Hex(), err)
	}
	return codeAt, nil
}

// CallContract calls the Ethereum RPC CallContract function.
func (d *Client) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.CallContract")
	var callResult []byte
	err := d.retryDo(ctx, func(ctx context.Context) error {
		var err error
		callResult, err = d.ethClient.CallContract(ctx, call, blockNumber)
		if err != nil {
			// NOTE: We have to do string comparison because generated bindings don't return a type we can easily work with
			// This error is not retryable
			if err.Error() == gethExecutionReverted || err.Error() == edgeExecutionReverted || err.Error() == hardhatExecutionReverted {
				return ErrExecutionReverted
			}
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return callResult, nil
}

func (d *Client) EthClient() *ethclient.Client {
	// TODO: can we avoid exposing this? Can we define an interface to accommodate?
	return d.ethClient
}

func (d *Client) retryDo(ctx context.Context, f retry.RetryFunc) error {
	// If no retries, just run the func
	if d.retryCount == 0 {
		return f(ctx)
	}
	// Run func with retries
	b := retry.WithMaxRetries(d.retryCount, retry.NewConstant(d.retryWait))
	err := retry.Do(ctx, b, f)
	return err
}

func logAndReturnRetryableErr(log *logrus.Entry, err error) error {
	log.Warnf("query error: %s", err.Error())
	return retry.RetryableError(err)
}
