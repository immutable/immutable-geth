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
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// The below functions implement the bind.ContractTransactor interface:
//
// - HeaderByNumber
// - PendingCodeAt
// - PendingNonceAt
// - SuggestGasPrice
// - SuggestGasTipCap
// - EstimateGas
// - SendTransaction
//
// These are used by bindings when executing contract calls. The reasoning
// behind having a custom implementation is that it allows for retry
// functionality and for custom gas pricing strategies.

func (client *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.HeaderByNumber")

	var header *types.Header
	err := client.retryDo(ctx, func(ctx context.Context) error {
		var err error
		header, err = client.ethClient.HeaderByNumber(ctx, number)
		if err != nil {
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})

	return header, err
}

func (client *Client) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.PendingCodeAt")

	var pendingCode []byte
	err := client.retryDo(ctx, func(ctx context.Context) error {
		var err error
		pendingCode, err = client.ethClient.PendingCodeAt(ctx, account)
		if err != nil {
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})

	return pendingCode, err
}

func (client *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.PendingNonceAt")

	var pendingNonce uint64
	err := client.retryDo(ctx, func(ctx context.Context) error {
		var err error
		pendingNonce, err = client.ethClient.PendingNonceAt(ctx, account)
		if err != nil {
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})

	return pendingNonce, err
}

// SuggestGasPrice is used for pricing legacy transactions.
//
// NOTE: Since we plan to use EIP-1559 transactions, there's no need to
// implement it. Panic to make sure anyone trying to use it knows it's not
// implemented.
func (client *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	panic("unimplemented")
}

func (client *Client) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.SuggestGasTipCap")
	log.Warnf("Temporary gas pricing calculation! Returns 0!")

	tip := big.NewInt(0)
	err := client.retryDo(ctx, func(ctx context.Context) error {
		var err error

		tip, err = client.gasTipEstimator.SuggestGasTipCap(ctx, client.ethClient)
		if err != nil {
			return logAndReturnRetryableErr(log, err)
		}

		return nil
	})

	return tip, err
}

func (client *Client) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	log := logrus.WithContext(ctx).WithField("function", "Client.EstimateGas")

	var gas uint64
	err := client.retryDo(ctx, func(ctx context.Context) error {
		var err error
		gas, err = client.ethClient.EstimateGas(ctx, call)
		if err != nil {
			// NOTE: We have to do string comparison because generated bindings
			// don't return a type we can easily work with This error is not
			// retryable.
			if err.Error() == gethExecutionReverted {
				return ErrExecutionReverted
			}
			return logAndReturnRetryableErr(log, err)
		}
		return nil
	})

	return gas, err
}

func (client *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	log := logrus.WithContext(ctx).WithField("function", "Client.SendTransaction")

	return client.retryDo(ctx, func(ctx context.Context) error {
		err := client.ethClient.SendTransaction(ctx, tx)
		if err != nil {
			return logAndReturnRetryableErr(log, err)
		}

		return nil
	})
}
